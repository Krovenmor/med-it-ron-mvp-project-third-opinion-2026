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

type apiClient struct {
	baseURL string
	http    *http.Client
}

func newAPIClient(baseURL string) *apiClient {
	return &apiClient{baseURL: baseURL, http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *apiClient) graph(t *testing.T, patientID string) (int, graphView) {
	t.Helper()
	status, body := c.do(t, http.MethodGet, "/api/v1/patients/"+patientID+"/graph", nil)
	if status != http.StatusOK {
		return status, graphView{}
	}
	return status, decode[graphView](t, body)
}

func (c *apiClient) slots(t *testing.T, patientID, recID string) (int, slotsView) {
	t.Helper()
	status, body := c.do(t, http.MethodGet, "/api/v1/patients/"+patientID+"/recommendations/"+recID+"/slots", nil)
	if status != http.StatusOK {
		return status, slotsView{}
	}
	return status, decode[slotsView](t, body)
}

func (c *apiClient) book(t *testing.T, patientID, recID string, req bookRequest) (int, appointmentView) {
	t.Helper()
	status, body := c.do(t, http.MethodPost, "/api/v1/patients/"+patientID+"/recommendations/"+recID+"/appointments", mustJSON(t, req))
	if status != http.StatusCreated {
		return status, appointmentView{}
	}
	return status, decode[appointmentView](t, body)
}

func (c *apiClient) do(t *testing.T, method, path string, body []byte) (int, []byte) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, respBody
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
