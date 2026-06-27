package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// errNotFound is returned by the client when the API responds 404, so resources
// can remove themselves from state instead of erroring.
var errNotFound = errors.New("bagre: resource not found")

// Client is a thin wrapper over the Bagre REST API authenticated with an API
// token (Authorization: Bearer bagre_…).
type Client struct {
	endpoint string
	token    string
	http     *http.Client
	ua       string
}

func NewClient(endpoint, token, version string) *Client {
	return &Client{
		endpoint: strings.TrimRight(endpoint, "/"),
		token:    token,
		http:     &http.Client{Timeout: 30 * time.Second},
		ua:       "terraform-provider-bagre/" + version,
	}
}

// do performs an authenticated JSON request. body and out may be nil.
func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", c.ua)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return errNotFound
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("bagre API %s %s: %d %s", method, path, resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if out != nil && len(data) > 0 {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("decoding %s %s response: %w", method, path, err)
		}
	}
	return nil
}

// --- Site entity (/api/sites) ---

type Site struct {
	ID          int64  `json:"id,omitempty"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (c *Client) CreateSite(ctx context.Context, s Site) (*Site, error) {
	var out Site
	if err := c.do(ctx, http.MethodPost, "/api/sites", s, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetSite(ctx context.Context, id int64) (*Site, error) {
	var out Site
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/sites/%d", id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateSite(ctx context.Context, id int64, s Site) (*Site, error) {
	var out Site
	if err := c.do(ctx, http.MethodPatch, fmt.Sprintf("/api/sites/%d", id), s, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteSite(ctx context.Context, id int64) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/sites/%d", id), nil, nil)
}
