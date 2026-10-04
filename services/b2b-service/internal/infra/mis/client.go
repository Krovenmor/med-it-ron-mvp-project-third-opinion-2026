package mis

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *Client) PatientHistory(ctx context.Context, p domain.Patient) (domain.PatientHistory, error) {
	endpoint := c.baseURL + "/api/v1/patients/" + url.PathEscape(p.ExternalID) + "/history"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return domain.PatientHistory{}, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.PatientHistory{}, fmt.Errorf("call mis: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return domain.PatientHistory{}, fmt.Errorf("mis responded %d: %s", resp.StatusCode, snippet)
	}

	var out historyResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return domain.PatientHistory{}, fmt.Errorf("decode response: %w", err)
	}
	return out.toDomain(), nil
}
