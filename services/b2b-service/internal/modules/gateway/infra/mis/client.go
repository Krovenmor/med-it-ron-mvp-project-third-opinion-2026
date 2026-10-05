package mis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/gateway/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/apperr"
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

func (c *Client) PatientHistory(ctx context.Context, externalID string) (domain.PatientHistory, error) {
	var out historyResponse
	status, err := c.do(ctx, http.MethodGet, "/api/v1/patients/"+url.PathEscape(externalID)+"/history", nil, &out)
	switch {
	case err != nil:
		return domain.PatientHistory{}, err
	case status != http.StatusOK:
		return domain.PatientHistory{}, unexpected(status)
	}
	return out.toDomain(), nil
}

func (c *Client) Slots(ctx context.Context, serviceCode string, from time.Time) ([]domain.Slot, error) {
	query := url.Values{"service_code": {serviceCode}, "from": {from.Format(time.RFC3339)}}
	var out slotsResponse
	status, err := c.do(ctx, http.MethodGet, "/api/v1/slots?"+query.Encode(), nil, &out)
	switch {
	case err != nil:
		return nil, err
	case status != http.StatusOK:
		return nil, unexpected(status)
	}
	return out.toDomain(), nil
}

func (c *Client) Book(ctx context.Context, req domain.AppointmentRequest) (domain.Appointment, error) {
	var out appointmentDTO
	status, err := c.do(ctx, http.MethodPost, "/api/v1/appointments", newAppointmentRequest(req), &out)
	switch {
	case err != nil:
		return domain.Appointment{}, err
	case status == http.StatusOK, status == http.StatusCreated:
		return out.toDomain(), nil
	case status == http.StatusConflict:
		return domain.Appointment{}, fmt.Errorf("%w: slot %s is already taken", apperr.ErrConflict, req.SlotID)
	case status == http.StatusBadRequest:
		return domain.Appointment{}, fmt.Errorf("%w: slot %s is not available for booking", apperr.ErrInvalidInput, req.SlotID)
	}
	return domain.Appointment{}, unexpected(status)
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) (int, error) {
	var payload io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return 0, fmt.Errorf("encode mis request: %w", err)
		}
		payload = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, payload)
	if err != nil {
		return 0, fmt.Errorf("build mis request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%w: call mis: %v", apperr.ErrUpstream, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusMultipleChoices {
		return resp.StatusCode, nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return 0, fmt.Errorf("%w: decode mis response: %v", apperr.ErrUpstream, err)
	}
	return resp.StatusCode, nil
}

func unexpected(status int) error {
	return fmt.Errorf("%w: mis responded %d", apperr.ErrUpstream, status)
}
