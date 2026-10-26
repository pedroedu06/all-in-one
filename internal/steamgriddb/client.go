// Package steamgriddb talks to the SteamGridDB API (cover art).
package steamgriddb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://www.steamgriddb.com/api/v2"

var (
	ErrInvalidKey  = errors.New("invalid SteamGridDB API key")
	ErrUnreachable = errors.New("SteamGridDB unreachable")
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New() *Client {
	return &Client{baseURL: defaultBaseURL, http: &http.Client{Timeout: 10 * time.Second}}
}

// ValidateKey makes a cheap authenticated request and reports whether the API
// accepts the key.
func (c *Client) ValidateKey(ctx context.Context, key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return ErrInvalidKey
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/search/autocomplete/test", nil)
	if err != nil {
		return fmt.Errorf("build SteamGridDB request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnreachable, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrInvalidKey
	default:
		return fmt.Errorf("SteamGridDB returned unexpected status %d", resp.StatusCode)
	}
}
