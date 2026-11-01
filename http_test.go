package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBearerAuth(t *testing.T) {
	var reached bool
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { reached = true })
	h := bearerAuth("secret", next)

	cases := []struct {
		header string
		want   int
		pass   bool
	}{
		{"", http.StatusUnauthorized, false},
		{"Bearer wrong", http.StatusUnauthorized, false},
		{"secret", http.StatusUnauthorized, false},
		{"Bearer secret", http.StatusOK, true},
	}
	for _, c := range cases {
		reached = false
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		if c.header != "" {
			req.Header.Set("Authorization", c.header)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if c.pass {
			if !reached {
				t.Errorf("header %q: handler not reached", c.header)
			}
		} else {
			if reached {
				t.Errorf("header %q: handler reached without valid token", c.header)
			}
			if rec.Code != c.want {
				t.Errorf("header %q: status %d, want %d", c.header, rec.Code, c.want)
			}
		}
	}
}
