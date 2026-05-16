package common

import (
	"crypto/x509"
	"errors"
	"fmt"
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
			want:      []string{"dns-01", "http-01"},
		},
		{
			name:      "dns-01 only config still falls back to http-01",
			preferred: "dns-01",
			enabled:   []string{"dns-01"},
			domains:   []string{"example.com"},
			want:      []string{"dns-01", "http-01"},
		},
		{
			name:      "wildcard dns-01 only does not gain http-01",
			preferred: "dns-01",
			enabled:   []string{"dns-01"},
			domains:   []string{"*.example.com"},
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

func TestDropFailedNames(t *testing.T) {
	cases := []struct {
		name        string
		domains     []string
		failed      []string
		wantRemain  []string
		wantDropped []string
	}{
		{
			name:        "drops a failing parked domain, keeps the primary",
			domains:     []string{"example.com", "parked.com", "www.example.com"},
			failed:      []string{"parked.com"},
			wantRemain:  []string{"example.com", "www.example.com"},
			wantDropped: []string{"parked.com"},
		},
		{
			name:        "never drops the primary even if it failed",
			domains:     []string{"example.com", "parked.com"},
			failed:      []string{"example.com", "parked.com"},
			wantRemain:  []string{"example.com"},
			wantDropped: []string{"parked.com"},
		},
		{
			name:        "nothing failed leaves the list intact",
			domains:     []string{"example.com", "parked.com"},
			failed:      nil,
			wantRemain:  []string{"example.com", "parked.com"},
			wantDropped: nil,
		},
		{
			name:        "drops several failing names, preserving order",
			domains:     []string{"example.com", "a.com", "b.com", "c.com"},
			failed:      []string{"c.com", "a.com"},
			wantRemain:  []string{"example.com", "b.com"},
			wantDropped: []string{"a.com", "c.com"},
		},
		{
			name:        "a failing wildcard name can be dropped",
			domains:     []string{"example.com", "*.example.com"},
			failed:      []string{"*.example.com"},
			wantRemain:  []string{"example.com"},
			wantDropped: []string{"*.example.com"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			remain, dropped := dropFailedNames(tc.domains, tc.failed)
			if !reflect.DeepEqual(remain, tc.wantRemain) {
				t.Errorf("remaining = %v, want %v", remain, tc.wantRemain)
			}
			if !reflect.DeepEqual(dropped, tc.wantDropped) {
				t.Errorf("dropped = %v, want %v", dropped, tc.wantDropped)
			}
		})
	}
}

func TestValidationFailedIdentifiers(t *testing.T) {
	ve := &ValidationError{Identifiers: []string{"parked.com"}, Err: errors.New("boom")}

	if got := ve.Error(); got != "boom" {
		t.Errorf("Error() = %q, want %q", got, "boom")
	}

	bare := &ValidationError{Identifiers: []string{"a.com", "b.com"}}
	if got := bare.Error(); got != "validation failed for: a.com, b.com" {
		t.Errorf("bare Error() = %q", got)
	}

	if got := validationFailedIdentifiers(ve); !reflect.DeepEqual(got, []string{"parked.com"}) {
		t.Errorf("direct: validationFailedIdentifiers = %v", got)
	}

	wrapped := fmt.Errorf("issuance failed: %w", ve)
	if got := validationFailedIdentifiers(wrapped); !reflect.DeepEqual(got, []string{"parked.com"}) {
		t.Errorf("wrapped: validationFailedIdentifiers = %v", got)
	}

	if got := validationFailedIdentifiers(errors.New("network down")); got != nil {
		t.Errorf("non-validation error: got %v, want nil", got)
	}
}
