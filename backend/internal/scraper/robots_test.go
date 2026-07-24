package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckRobotsAllowed(t *testing.T) {
	tests := []struct {
		name        string
		robotsBody  string
		path        string
		wantAllowed bool
	}{
		{
			name:        "allow all",
			robotsBody:  "User-agent: *\nAllow: /",
			path:        "/collections/all-products",
			wantAllowed: true,
		},
		{
			name:        "disallow specific path",
			robotsBody:  "User-agent: *\nDisallow: /collections/all-products",
			path:        "/collections/all-products",
			wantAllowed: false,
		},
		{
			name:        "disallow everything",
			robotsBody:  "User-agent: *\nDisallow: /",
			path:        "/collections/all-products",
			wantAllowed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tt.robotsBody))
			}))
			defer server.Close()

			err := CheckRobotsAllowed(server.URL+tt.path, "KickPickBot/1.0")
			allowed := err == nil
			if allowed != tt.wantAllowed {
				t.Errorf("CheckRobotsAllowed() allowed = %v, want %v (err: %v)", allowed, tt.wantAllowed, err)
			}
		})
	}
}

func TestCheckRobotsAllowed_FailsClosedOnFetchError(t *testing.T) {
	// Nothing listening on this address — the fetch itself should fail, and
	// per the fail-closed contract that must be treated as "not allowed".
	err := CheckRobotsAllowed("http://127.0.0.1:1/some-path", "KickPickBot/1.0")
	if err == nil {
		t.Error("expected an error when robots.txt can't be fetched, got nil")
	}
}
