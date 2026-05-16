package common

import (
	"errors"
	"testing"
)

func TestIsWhmAuthError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"WHM 401 status line", errors.New("401 Access Denied"), true},
		{"403 status line", errors.New("403 Forbidden"), true},
		{"access denied inside a wrapped message", errors.New("fetching vhosts: Access Denied"), true},
		{"unauthorized", errors.New("401 Unauthorized"), true},
		{"connection refused is not recoverable", errors.New("dial tcp 127.0.0.1:2087: connect: connection refused"), false},
		{"server error is not an auth error", errors.New("500 Internal Server Error"), false},
		{"unrelated api error", errors.New("API response maximum size exceeded"), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsWhmAuthError(tc.err); got != tc.want {
				t.Errorf("IsWhmAuthError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
