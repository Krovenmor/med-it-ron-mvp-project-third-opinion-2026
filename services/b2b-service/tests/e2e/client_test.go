//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	maxAttempts    = 3
	settleTimeout  = 15 * time.Second
	settleInterval = 100 * time.Millisecond
)

type apiClient struct {
	baseURL string
	http    *http.Client
}

func newAPIClient(baseURL string) *apiClient {
	return &apiClient{baseURL: baseURL, http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *apiClient) ingest(t *testing.T, r report) (int, caseRef) {
	t.Helper()
	status, body := c.do(t, http.MethodPost, "/api/v1/reports", mustJSON(t, r))
	if status != http.StatusOK && status != http.StatusCreated {
		return status, caseRef{}
	}
	return status, decode[caseRef](t, body)
}

func (c *apiClient) getCase(t *testing.T, id string) (int, caseView) {
	t.Helper()
	status, body := c.do(t, http.MethodGet, "/api/v1/cases/"+id, nil)
	if status != http.StatusOK {
		return status, caseView{}
	}
	return status, decode[caseView](t, body)
}

func (c *apiClient) awaitStatus(t *testing.T, id, want string) caseView {
	t.Helper()
	var (
		status int
		last   caseView
	)
	reached := poll(func() bool {
		status, last = c.getCase(t, id)
		return status == http.StatusOK && last.Status == want
	})
	require.True(t, reached, "case %s: http %d, status %q, want %q", id, status, last.Status, want)
	return last
}

func (c *apiClient) reviewQueue(t *testing.T) reviewQueueView {
	t.Helper()
	status, body := c.do(t, http.MethodGet, "/api/v1/review/queue", nil)
	require.Equal(t, http.StatusOK, status, "body: %s", body)
	return decode[reviewQueueView](t, body)
}

func (c *apiClient) openCase(t *testing.T, actor, caseID string) int {
	t.Helper()
	status, _ := c.doAs(t, actor, http.MethodPost, "/api/v1/cases/"+caseID+"/open", nil)
	return status
}

func (c *apiClient) reviewRecommendation(t *testing.T, actor, caseID, recID string, req reviewRequest) (int, caseRecommendation) {
	t.Helper()
	status, body := c.doAs(t, actor, http.MethodPatch, "/api/v1/cases/"+caseID+"/recommendations/"+recID, mustJSON(t, req))
	if status != http.StatusOK {
		return status, caseRecommendation{}
	}
	return status, decode[caseRecommendation](t, body)
}

func (c *apiClient) addRecommendation(t *testing.T, actor, caseID string, req addRecommendationRequest) (int, caseRecommendation) {
	t.Helper()
	status, body := c.doAs(t, actor, http.MethodPost, "/api/v1/cases/"+caseID+"/recommendations", mustJSON(t, req))
	if status != http.StatusCreated {
		return status, caseRecommendation{}
	}
	return status, decode[caseRecommendation](t, body)
}

func (c *apiClient) changeUrgency(t *testing.T, actor, caseID string, req urgencyRequest) int {
	t.Helper()
	status, _ := c.doAs(t, actor, http.MethodPut, "/api/v1/cases/"+caseID+"/urgency", mustJSON(t, req))
	return status
}

func (c *apiClient) confirm(t *testing.T, actor, caseID string) (int, caseRef) {
	t.Helper()
	status, body := c.doAs(t, actor, http.MethodPost, "/api/v1/cases/"+caseID+"/confirm", nil)
	if status != http.StatusOK {
		return status, caseRef{}
	}
	return status, decode[caseRef](t, body)
}

func poll(cond func() bool) bool {
	deadline := time.Now().Add(settleTimeout)
	for !cond() {
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(settleInterval)
	}
	return true
}

func (c *apiClient) do(t *testing.T, method, path string, body []byte) (int, []byte) {
	t.Helper()
	return c.doAs(t, "", method, path, body)
}

func (c *apiClient) doAs(t *testing.T, actor, method, path string, body []byte) (int, []byte) {
	t.Helper()
	status, respBody, err := c.send(method, path, body, actor)
	require.NoError(t, err)
	return status, respBody
}

func (c *apiClient) send(method, path string, body []byte, actor string) (int, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if actor != "" {
		req.Header.Set("X-User-ID", actor)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, respBody, nil
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	require.NoError(t, err)
	return data
}

func decode[T any](t *testing.T, data []byte) T {
	t.Helper()
	var v T
	require.NoError(t, json.Unmarshal(data, &v), "body: %s", data)
	return v
}
