package common

import (
	"strings"
	"time"

	"github.com/domainr/dnsr"
	"github.com/letsencrypt-cpanel/cpanelgo/cpanel"
)

// dnsViabilityTimeout bounds the public DNS lookup so the issuance page never
// hangs on a slow or unreachable resolver.
const dnsViabilityTimeout = 5 * time.Second

// IsDNS01Viable reports whether dns-01 validation, and therefore wildcard
// issuance, is likely to succeed for domain.
//
// It is true only when cPanel hosts an authoritative DNS zone for the domain
// and the domain's public NS delegation points back at the nameservers cPanel
// serves for that zone. When the registrar delegates DNS elsewhere (Cloudflare,
// the registrar's own DNS and so on) the TXT record cPanel writes is never seen
// by Let's Encrypt, so a wildcard cannot be validated.
//
// It performs network I/O (cPanel API calls and an iterative DNS lookup) and is
// intended for UI hints: any failure is reported as "not viable" so the caller
// falls back to the safe default of not pre-selecting a wildcard.
func IsDNS01Viable(cp cpanel.CpanelApi, domain string) bool {
	domain = normalizeDNSName(domain)
	if domain == "" {
		return false
	}

	// cPanel must host a DNS zone for the domain.
	zones, err := cp.FetchZones()
	if err != nil {
		return false
	}
	zoneRoot := zones.FindRootForName(domain)
	if zoneRoot == "" {
		return false
	}

	// The nameservers cPanel publishes for that zone must match the domain's
	// real, public delegation, otherwise cPanel is not authoritative for it.
	zone, err := cp.FetchZone(zoneRoot, "NS")
	if err != nil || len(zone.Data) == 0 {
		return false
	}
	var localNS []string
	for _, rr := range zone.Data[0].Records {
		if strings.EqualFold(rr.Type, "NS") {
			localNS = append(localNS, rr.Record)
		}
	}

	var publicNS []string
	for _, rr := range dnsr.NewWithTimeout(0, dnsViabilityTimeout).Resolve(zoneRoot, "NS") {
		if strings.EqualFold(rr.Type, "NS") {
			publicNS = append(publicNS, rr.Value)
		}
	}

	return nameserverSetsOverlap(localNS, publicNS)
}

// normalizeDNSName lowercases a DNS name and strips a trailing dot so that
// values from different sources (cPanel zone records and wire-format DNS
// responses) can be compared.
func normalizeDNSName(name string) string {
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(name)), ".")
}

// nameserverSetsOverlap reports whether two nameserver lists share at least one
// host. It is the pure decision behind IsDNS01Viable: when the public
// delegation and the cPanel-hosted zone agree on a nameserver, cPanel is
// authoritative for the domain and dns-01 can be solved.
func nameserverSetsOverlap(local, public []string) bool {
	hosts := map[string]bool{}
	for _, ns := range local {
		if n := normalizeDNSName(ns); n != "" {
			hosts[n] = true
		}
	}
	for _, ns := range public {
		if n := normalizeDNSName(ns); n != "" && hosts[n] {
			return true
		}
	}
	return false
}
