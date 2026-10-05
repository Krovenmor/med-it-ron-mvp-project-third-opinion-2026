package aiservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
)

type Client struct {
	endpoint string
	http     *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		endpoint: strings.TrimRight(baseURL, "/") + "/v1/assessments",
		http:     &http.Client{Timeout: timeout},
	}
}

func (c *Client) Assess(ctx context.Context, req domain.AssessmentRequest) (domain.Assessment, error) {
	body, err := json.Marshal(newAssessmentRequest(req))
	if err != nil {
		return domain.Assessment{}, fmt.Errorf("encode request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return domain.Assessment{}, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return domain.Assessment{}, fmt.Errorf("call ai-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return domain.Assessment{}, fmt.Errorf("ai-service responded %d: %s", resp.StatusCode, snippet)
	}

	var out assessmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return domain.Assessment{}, fmt.Errorf("decode response: %w", err)
	}
	return out.toDomain(), nil
}
