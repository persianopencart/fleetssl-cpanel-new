package common

import "testing"

func TestNormalizeDNSName(t *testing.T) {
	cases := map[string]string{
		"NS1.Example.COM.":     "ns1.example.com",
		"ns1.example.com":      "ns1.example.com",
		"  ns1.example.com.  ": "ns1.example.com",
		"":                     "",
		".":                    "",
	}
	for in, want := range cases {
		if got := normalizeDNSName(in); got != want {
			t.Errorf("normalizeDNSName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNameserverSetsOverlap(t *testing.T) {
	tests := []struct {
		name   string
		local  []string
		public []string
		want   bool
	}{
		{
			name:   "registrar delegates to the cPanel server",
			local:  []string{"ns1.hostingco.com.", "ns2.hostingco.com."},
			public: []string{"ns1.hostingco.com.", "ns2.hostingco.com."},
			want:   true,
		},
		{
			name:   "case and trailing-dot differences still match",
			local:  []string{"NS1.HostingCo.com"},
			public: []string{"ns1.hostingco.com."},
			want:   true,
		},
		{
			name:   "domain delegated to external DNS",
			local:  []string{"ns1.hostingco.com."},
			public: []string{"gail.ns.cloudflare.com.", "rick.ns.cloudflare.com."},
			want:   false,
		},
		{
			name:   "a single shared nameserver is enough",
			local:  []string{"ns1.hostingco.com.", "ns2.hostingco.com."},
			public: []string{"ns1.hostingco.com.", "ns2.otherdns.net."},
			want:   true,
		},
		{
			name:   "no public delegation resolved",
			local:  []string{"ns1.hostingco.com."},
			public: nil,
			want:   false,
		},
		{
			name:   "empty or root names never match",
			local:  []string{"", "."},
			public: []string{"", "."},
			want:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := nameserverSetsOverlap(tc.local, tc.public); got != tc.want {
				t.Errorf("nameserverSetsOverlap(%v, %v) = %v, want %v",
					tc.local, tc.public, got, tc.want)
			}
		})
	}
}
