package steamgriddb

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := New()
	c.baseURL = srv.URL
	return c
}

// SteamGridDB authenticates with a Bearer token; sending it any other way
// makes every valid key look invalid.
func TestValidateKeySendsBearerToken(t *testing.T) {
	var got string
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
	})

	if err := c.ValidateKey(context.Background(), "  abc123  "); err != nil {
		t.Fatalf("ValidateKey: %v", err)
	}
	if got != "Bearer abc123" {
		t.Fatalf("Authorization = %q, want %q", got, "Bearer abc123")
	}
}

// The settings screen distinguishes a wrong key from other failures so the
// user knows whether to fix the key or retry later.
func TestValidateKeyClassifiesFailures(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		wantErr error
	}{
		{"unauthorized means bad key", http.StatusUnauthorized, ErrInvalidKey},
		{"forbidden means bad key", http.StatusForbidden, ErrInvalidKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(tt.status) })
			if err := c.ValidateKey(context.Background(), "k"); !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}

	t.Run("server error is not reported as bad key", func(t *testing.T) {
		c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) })
		err := c.ValidateKey(context.Background(), "k")
		if err == nil || errors.Is(err, ErrInvalidKey) || !strings.Contains(err.Error(), "500") {
			t.Fatalf("err = %v, want unexpected-status error mentioning 500", err)
		}
	})
}

func TestValidateKeyReportsUnreachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	srv.Close()
	c := New()
	c.baseURL = srv.URL

	if err := c.ValidateKey(context.Background(), "k"); !errors.Is(err, ErrUnreachable) {
		t.Fatalf("err = %v, want ErrUnreachable", err)
	}
}

// An empty key is obviously invalid; don't spend a network round-trip on it.
func TestValidateKeyEmptyKeySkipsRequest(t *testing.T) {
	called := false
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { called = true })

	if err := c.ValidateKey(context.Background(), "   "); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("err = %v, want ErrInvalidKey", err)
	}
	if called {
		t.Fatal("request was made for an empty key")
	}
}
