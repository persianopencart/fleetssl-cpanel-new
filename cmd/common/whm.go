package common

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/letsencrypt-cpanel/cpanelgo/whm"
	log "github.com/sirupsen/logrus"
)

// MakeWhmClient builds a WHM API client authenticated with the stored API
// token. The token is validated before the client is returned; if WHM rejects
// it as stale, the token is discarded, regenerated, and the client is rebuilt
// once — so a revoked or carried-over token recovers automatically instead of
// failing every request until the token file is removed by hand.
func MakeWhmClient(insecure bool) (whm.WhmApi, error) {
	cl, err := newWhmClient(insecure)
	if err != nil {
		return whm.WhmApi{}, err
	}

	if _, err := cl.Version(); err != nil && IsWhmAuthError(err) {
		log.WithError(err).Warn("WHM rejected the stored API token; regenerating it")
		if ierr := InvalidateApiToken(); ierr != nil {
			log.WithError(ierr).Error("Could not discard the stale WHM API token")
			return cl, nil
		}
		return newWhmClient(insecure)
	}

	return cl, nil
}

func newWhmClient(insecure bool) (whm.WhmApi, error) {
	s, err := ReadApiToken()
	if err != nil {
		return whm.WhmApi{}, err
	}

	hn, err := os.Hostname()
	if err != nil {
		return whm.WhmApi{}, err
	}

	return whm.NewWhmApiAccessHashTotp(hn, "root", s, insecure, ReadTotpSecret()), nil
}

// IsWhmAuthError reports whether err is a WHM API authentication failure — a
// rejected or expired access token. These are recoverable by regenerating the
// token; network and server-side errors are not, and must not be misclassified.
func IsWhmAuthError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "access denied") ||
		strings.Contains(msg, "unauthorized") ||
		strings.Contains(msg, "permission denied") ||
		strings.HasPrefix(msg, "401 ") ||
		strings.HasPrefix(msg, "403 ")
}

func ReadTotpSecret() string {
	totpBytes, err := os.ReadFile("/var/cpanel/authn/twofactor_auth/tfa_userdata.json")
	if err != nil {
		return ""
	}

	var totpMap map[string]map[string]string

	err = json.Unmarshal(totpBytes, &totpMap)
	if err != nil {
		return ""
	}

	root, ok := totpMap["root"]
	if !ok {
		return ""
	}

	secret, ok := root["secret"]
	if !ok {
		return ""
	}

	return secret
}
