package httpjson

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
)

type Client struct {
	http *http.Client
}

func New(timeout time.Duration) Client {
	return Client{http: &http.Client{Timeout: timeout}}
}

func (c Client) Do(ctx context.Context, method, url string, in, out any) (int, error) {
	var body io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return 0, fmt.Errorf("encode request: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return 0, fmt.Errorf("build request: %w", err)
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%w: %s %s: %v", domain.ErrUpstream, method, url, err)
	}
	defer resp.Body.Close()

	if out == nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		_, _ = io.Copy(io.Discard, resp.Body)
		return resp.StatusCode, nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return resp.StatusCode, fmt.Errorf("%w: decode %s %s response: %v", domain.ErrUpstream, method, url, err)
	}
	return resp.StatusCode, nil
}

func Unexpected(service string, status int) error {
	return fmt.Errorf("%w: %s responded %d", domain.ErrUpstream, service, status)
}
