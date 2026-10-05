package b2b

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/demo"
)

type Client struct {
	endpoint string
	http     *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		endpoint: strings.TrimRight(baseURL, "/") + "/api/v1/reports",
		http:     &http.Client{Timeout: timeout},
	}
}

func (c *Client) SendReport(ctx context.Context, r demo.Report) (Delivery, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return Delivery{}, fmt.Errorf("encode report: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return Delivery{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return Delivery{}, fmt.Errorf("send report %s: %w", r.Study.ID, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return Delivery{}, fmt.Errorf("read response for %s: %w", r.Study.ID, err)
	}
	if !json.Valid(respBody) {
		respBody, _ = json.Marshal(string(respBody))
	}
	return Delivery{StudyID: r.Study.ID, Status: resp.StatusCode, Response: respBody}, nil
}
