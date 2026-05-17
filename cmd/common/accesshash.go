package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/letsencrypt-cpanel/cpanelgo/whm"
	log "github.com/sirupsen/logrus"
)

const (
	pathApiToken   = "/etc/.letsencrypt-cpanel-api-token"
	pathAccessHash = "/root/.accesshash"
	pathWhmapi1    = "/usr/local/cpanel/bin/whmapi1"

	// apiTokenNamePrefix is the name prefix of every WHM API token this plugin
	// creates. Tokens with this prefix are treated as the plugin's own and are
	// revoked when a fresh token is generated, so they cannot accumulate
	// against WHM's per-account API token limit.
	apiTokenNamePrefix = "letsencrypt-cpanel-autogen_"
)

var (
	apiTokenFailing = false
)

// ReadApiToken returns a WHM API token for authenticating as root. It reads the
// stored token (creating one if necessary) and only falls back to a legacy
// /root/.accesshash if token handling is unavailable.
func ReadApiToken() (string, error) {
	if !apiTokenFailing {
		token, err := getOrCreateApiToken()
		if err == nil {
			return token, nil
		}
		log.WithError(err).Error("Getting/creating WHM API token, attempting to fall back to accesshash")
		apiTokenFailing = true
	}

	if FileExists(pathAccessHash) {
		t, err := readFile(pathAccessHash)
		if err == nil {
			if t != "" {
				return t, nil
			}
			log.Error("access hash empty")
		} else {
			log.WithError(err).Error("reading access hash")
		}
	}

	return "", errors.New("No cPanel API authentication token available. " +
		"WHM may have reached its API token limit — revoke unused API tokens in " +
		"WHM (\"Manage API Tokens\") and try again.")
}

// getOrCreateApiToken returns the stored WHM API token, generating a fresh one
// whenever the token file is missing, empty or otherwise unreadable. The empty
// check matters: a zero-byte token file (left behind by an interrupted create)
// would otherwise wedge every subsequent request, because the old code only
// created a token when the file was entirely absent.
func getOrCreateApiToken() (string, error) {
	if token, err := readFile(pathApiToken); err == nil {
		return token, nil
	}

	if err := createApiToken(); err != nil {
		return "", err
	}

	return readFile(pathApiToken)
}

// InvalidateApiToken discards the stored WHM API token so that the next call to
// ReadApiToken generates a fresh one. It is used to recover from a token that
// WHM has begun rejecting — for example after the token was revoked, or the
// server was restored from a backup so the saved token no longer matches WHM's
// token store.
func InvalidateApiToken() error {
	apiTokenFailing = false
	if err := os.Remove(pathApiToken); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func readFile(path string) (string, error) {
	ahBytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Error reading api token file %s, %v", path, err)
	}
	if len(strings.TrimSpace(string(ahBytes))) == 0 {
		return "", fmt.Errorf("Api token file %s is empty!", path)
	}

	return string(ahBytes), nil
}

func getTokenName() string {
	return apiTokenNamePrefix + time.Now().Format("2006-01-02_15-04-05")
}

// defaultApiTokenAcls is the least-privilege set of WHM ACLs the plugin needs:
// impersonating cPanel users (for renewals, dns-01 and certificate installs),
// SSL inspection and installation, enumerating accounts, and restarting
// services. It deliberately excludes account creation and termination,
// password changes, DNS zone management, packages, reseller privileges and
// root access — so a leaked token cannot be used to take over the server.
var defaultApiTokenAcls = []string{
	"cpanel-api",          // run cPanel/UAPI calls as each account
	"ssl",                 // install and manage SSL, including service certificates
	"ssl-info",            // read the installed SSL virtual hosts
	"list-accts",          // enumerate accounts for renewal
	"acct-summary",        // per-account details during renewal
	"restart",             // restart cpsrvd/ftpd after a service-certificate install
	"basic-whm-functions", // the WHM version check and other basic calls
}

// pathTokenAcls is an optional administrator override: a comma- or
// newline-separated list of WHM ACLs to request for the API token. It is an
// escape hatch in case a cPanel version maps an API call to a different ACL
// than defaultApiTokenAcls expects, so the set can be corrected without
// rebuilding the plugin. Lines beginning with "#" are treated as comments.
const pathTokenAcls = "/etc/letsencrypt-cpanel-token-acls"

// apiTokenAcls returns the WHM ACLs to request for a new API token: the
// administrator override from pathTokenAcls when present and non-empty,
// otherwise the built-in least-privilege default.
func apiTokenAcls() []string {
	raw, err := os.ReadFile(pathTokenAcls)
	if err != nil {
		return defaultApiTokenAcls
	}
	var acls []string
	for _, line := range strings.Split(string(raw), "\n") {
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = line[:i]
		}
		for _, tok := range strings.Split(line, ",") {
			if tok = strings.TrimSpace(tok); tok != "" {
				acls = append(acls, tok)
			}
		}
	}
	if len(acls) == 0 {
		return defaultApiTokenAcls
	}
	return acls
}

// createApiToken generates a new WHM API token, scoped to the ACLs from
// apiTokenAcls so it carries only the privileges the plugin needs, and stores
// it. Tokens previously created by this plugin are revoked first: WHM caps how
// many API tokens may exist, and without this cleanup the plugin leaves a stale
// token behind on every install, reinstall and regeneration until token
// creation eventually fails with "No cPanel API authentication token available".
func createApiToken() error {
	revokeStaleApiTokens()

	tn := getTokenName()
	acls := apiTokenAcls()
	log.WithFields(log.Fields{"token_name": tn, "acls": acls}).Info("Creating new WHM API token")

	// The token must be created with ACLs: one created with none has no
	// privileges, so every WHM API call would fail with "Access Denied" and
	// the client would pointlessly regenerate the token on every request.
	args := []string{"--output=json", "api_token_create", "token_name=" + tn}
	for i, acl := range acls {
		args = append(args, fmt.Sprintf("acl-%d=%s", i+1, acl))
	}

	out, err := exec.Command(pathWhmapi1, args...).Output()
	if err != nil {
		return fmt.Errorf("Error calling new api token command: %v%s", err, exitErrorDetail(err))
	}

	result := struct {
		whm.BaseWhmApiResponse
		Data struct {
			Name       string `json:"name"`
			Token      string `json:"token"`
			CreateTime int    `json:"create_time"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(out, &result); err != nil {
		return fmt.Errorf("Error parsing create token output: %v", err)
	}
	if err := result.Error(); err != nil {
		return fmt.Errorf("Error creating new api token: %v", err)
	}
	if result.Data.Token == "" {
		return errors.New("WHM returned an empty API token from api_token_create")
	}

	if err := os.WriteFile(pathApiToken, []byte(result.Data.Token), ConfigPermissions); err != nil {
		return fmt.Errorf("Error saving new api token: %v", err)
	}

	return nil
}

// revokeStaleApiTokens revokes every WHM API token previously auto-created by
// this plugin (those named with apiTokenNamePrefix). It is best-effort: any
// failure is logged and ignored, because its only purpose is housekeeping so
// old tokens cannot accumulate against WHM's API token limit.
func revokeStaleApiTokens() {
	out, err := exec.Command(pathWhmapi1, "--output=json", "api_token_list").Output()
	if err != nil {
		log.WithError(err).Warn("Could not list WHM API tokens to clean up old ones")
		return
	}

	result := struct {
		whm.BaseWhmApiResponse
		Data struct {
			Tokens map[string]json.RawMessage `json:"tokens"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(out, &result); err != nil {
		log.WithError(err).Warn("Could not parse the WHM API token list")
		return
	}
	if err := result.Error(); err != nil {
		log.WithError(err).Warn("WHM rejected the api_token_list request")
		return
	}

	for name := range result.Data.Tokens {
		if !strings.HasPrefix(name, apiTokenNamePrefix) {
			continue
		}
		if _, err := exec.Command(pathWhmapi1, "--output=json", "api_token_revoke", "token_name="+name).Output(); err != nil {
			log.WithError(err).WithField("token_name", name).
				Warn("Failed to revoke a stale WHM API token")
			continue
		}
		log.WithField("token_name", name).Info("Revoked stale WHM API token")
	}
}

// exitErrorDetail returns the stderr captured from a failed exec.Command so the
// underlying whmapi1 error is visible instead of a bare "exit status N".
func exitErrorDetail(err error) string {
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		if s := strings.TrimSpace(string(ee.Stderr)); s != "" {
			return ": " + s
		}
	}
	return ""
}
