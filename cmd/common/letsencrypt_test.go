package common

import (
	"crypto/x509"
	"reflect"
	"testing"
)

func parseDERs(chain [][]byte) []*x509.Certificate {
	out := []*x509.Certificate{}
	for _, der := range chain {
		crt, err := x509.ParseCertificate(der)
		if err != nil {
			panic(err)
		}
		out = append(out, crt)
	}
	return out
}

func TestChallengeMethodChain(t *testing.T) {
	cases := []struct {
		name      string
		preferred string
		enabled   []string
		domains   []string
		want      []string
	}{
		{
			name:      "preferred goes first then remaining enabled",
			preferred: "http-01",
			enabled:   []string{"dns-01", "http-01"},
			domains:   []string{"example.com"},
			want:      []string{"http-01", "dns-01"},
		},
		{
			name:      "no preferred uses configured order",
			preferred: "",
			enabled:   []string{"dns-01", "http-01"},
			domains:   []string{"example.com"},
			want:      []string{"dns-01", "http-01"},
		},
		{
			name:      "whitespace trimmed and duplicates dropped",
			preferred: " dns-01 ",
			enabled:   []string{"dns-01", " http-01 "},
			domains:   []string{"example.com"},
			want:      []string{"dns-01", "http-01"},
		},
		{
			name:      "invalid methods are ignored",
			preferred: "tls-alpn-01",
			enabled:   []string{"bogus", "dns-01"},
			domains:   []string{"example.com"},
			want:      []string{"dns-01"},
		},
		{
			name:      "wildcard drops http-01 from the chain",
			preferred: "http-01",
			enabled:   []string{"http-01", "dns-01"},
			domains:   []string{"*.example.com", "example.com"},
			want:      []string{"dns-01"},
		},
		{
			name:      "wildcard with no dns-01 still falls back to dns-01",
			preferred: "http-01",
			enabled:   []string{"http-01"},
			domains:   []string{"*.example.com"},
			want:      []string{"dns-01"},
		},
		{
			name:      "empty configuration defaults to http-01",
			preferred: "",
			enabled:   nil,
			domains:   []string{"example.com"},
			want:      []string{"http-01"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ChallengeMethodChain(tc.preferred, tc.enabled, tc.domains)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ChallengeMethodChain(%q, %v, %v) = %v, want %v",
					tc.preferred, tc.enabled, tc.domains, got, tc.want)
			}
		})
	}
}
