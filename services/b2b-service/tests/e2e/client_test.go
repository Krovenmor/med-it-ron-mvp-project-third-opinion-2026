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

func (c *apiClient) reportStatus(t *testing.T, caseID string) (int, reportStatusView) {
	t.Helper()
	status, body := c.do(t, http.MethodGet, "/api/v1/reports/"+caseID, nil)
	if status != http.StatusOK {
		return status, reportStatusView{}
	}
	return status, decode[reportStatusView](t, body)
}

func (c *apiClient) routeStatus(t *testing.T, cv caseView) string {
	t.Helper()
	status, plan, _ := c.plan(t, cv.SourceSystem, cv.Patient.ID)
	if status != http.StatusOK {
		return ""
	}
	for _, pc := range plan.Cases {
		if pc.CaseID == cv.ID {
			return pc.Status
		}
	}
	return ""
}

func (c *apiClient) awaitRouteStatus(t *testing.T, cv caseView, want string) {
	t.Helper()
	var last string
	reached := poll(func() bool {
		last = c.routeStatus(t, cv)
		return last == want
	})
	require.True(t, reached, "case %s: route status %q, want %q", cv.ID, last, want)
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

func (c *apiClient) plan(t *testing.T, sourceSystem, patientID string) (int, planView, []byte) {
	t.Helper()
	status, body := c.do(t, http.MethodGet, "/api/v1/patients/"+patientID+"/plan?source_system="+sourceSystem, nil)
	if status != http.StatusOK {
		return status, planView{}, body
	}
	return status, decode[planView](t, body), body
}

func (c *apiClient) registerBooking(t *testing.T, caseID string, req bookingRequest) (int, bookingView) {
	t.Helper()
	status, body := c.do(t, http.MethodPost, "/api/v1/cases/"+caseID+"/bookings", mustJSON(t, req))
	if status != http.StatusOK && status != http.StatusCreated {
		return status, bookingView{}
	}
	return status, decode[bookingView](t, body)
}

func (c *apiClient) declineRecommendation(t *testing.T, caseID, recID string, req patientDeclineRequest) (int, patientDeclineView) {
	t.Helper()
	status, body := c.do(t, http.MethodPost, "/api/v1/cases/"+caseID+"/recommendations/"+recID+"/decline", mustJSON(t, req))
	if status != http.StatusOK {
		return status, patientDeclineView{}
	}
	return status, decode[patientDeclineView](t, body)
}

func (c *apiClient) requestHelp(t *testing.T, caseID string) (int, helpRequestView) {
	t.Helper()
	status, body := c.do(t, http.MethodPost, "/api/v1/cases/"+caseID+"/help-requests", nil)
	if status != http.StatusOK && status != http.StatusCreated {
		return status, helpRequestView{}
	}
	return status, decode[helpRequestView](t, body)
}

func (c *apiClient) caseHistory(t *testing.T, caseID string) (int, history) {
	t.Helper()
	status, body := c.do(t, http.MethodGet, "/api/v1/cases/"+caseID+"/history", nil)
	if status != http.StatusOK {
		return status, history{}
	}
	return status, decode[history](t, body)
}

func (c *apiClient) clock(t *testing.T) time.Time {
	t.Helper()
	status, body := c.do(t, http.MethodGet, "/api/v1/clock", nil)
	require.Equal(t, http.StatusOK, status)
	return decode[clockView](t, body).Now
}

func (c *apiClient) advanceClock(t *testing.T, by time.Duration) time.Time {
	t.Helper()
	status, body := c.do(t, http.MethodPost, "/demo/clock/advance", mustJSON(t, advanceClockRequest{By: by.String()}))
	require.Equal(t, http.StatusOK, status, "body: %s", body)
	return decode[clockView](t, body).Now
}

func (c *apiClient) operatorTasks(t *testing.T, actor string, mine bool) operatorTasksView {
	t.Helper()
	path := "/api/v1/operator/tasks"
	if mine {
		path += "?mine=true"
	}
	status, body := c.doAs(t, actor, http.MethodGet, path, nil)
	require.Equal(t, http.StatusOK, status, "body: %s", body)
	return decode[operatorTasksView](t, body)
}

func (c *apiClient) operatorCard(t *testing.T, taskID string) (int, operatorCardView) {
	t.Helper()
	status, body := c.do(t, http.MethodGet, "/api/v1/operator/tasks/"+taskID, nil)
	if status != http.StatusOK {
		return status, operatorCardView{}
	}
	return status, decode[operatorCardView](t, body)
}

func (c *apiClient) operatorAction(t *testing.T, actor, taskID, action string, req any) (int, operatorTaskView) {
	t.Helper()
	var payload []byte
	if req != nil {
		payload = mustJSON(t, req)
	}
	status, body := c.doAs(t, actor, http.MethodPost, "/api/v1/operator/tasks/"+taskID+"/"+action, payload)
	if status != http.StatusOK {
		return status, operatorTaskView{}
	}
	return status, decode[operatorTaskView](t, body)
}

func (c *apiClient) operatorBook(t *testing.T, actor, taskID string, req operatorBookingRequest) (int, bookingView) {
	t.Helper()
	status, body := c.doAs(t, actor, http.MethodPost, "/api/v1/operator/tasks/"+taskID+"/bookings", mustJSON(t, req))
	if status != http.StatusOK && status != http.StatusCreated {
		return status, bookingView{}
	}
	return status, decode[bookingView](t, body)
}

func (c *apiClient) notifications(t *testing.T) []notificationView {
	t.Helper()
	status, body := c.do(t, http.MethodGet, "/api/v1/notifications", nil)
	require.Equal(t, http.StatusOK, status, "body: %s", body)
	return decode[notificationsView](t, body).Notifications
}

func (c *apiClient) dashboard(t *testing.T, days string) (int, dashboardView) {
	t.Helper()
	status, body := c.do(t, http.MethodGet, "/api/v1/dashboard?days="+days, nil)
	if status != http.StatusOK {
		return status, dashboardView{}
	}
	return status, decode[dashboardView](t, body)
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
