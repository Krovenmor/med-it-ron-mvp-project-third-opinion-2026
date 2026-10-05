package mis

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/domain"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/infra/httpjson"
)

const service = "mis"

type Client struct {
	baseURL string
	http    httpjson.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: httpjson.New(timeout)}
}

func (c *Client) History(ctx context.Context, patientID string) (domain.History, error) {
	var out historyResponse
	status, err := c.http.Do(ctx, http.MethodGet, c.baseURL+"/api/v1/patients/"+url.PathEscape(patientID)+"/history", nil, &out)
	switch {
	case err != nil:
		return domain.History{}, err
	case status != http.StatusOK:
		return domain.History{}, httpjson.Unexpected(service, status)
	}
	return out.toDomain(), nil
}

func (c *Client) Slots(ctx context.Context, serviceCode string) ([]domain.Slot, error) {
	var out slotsResponse
	status, err := c.http.Do(ctx, http.MethodGet, c.baseURL+"/api/v1/slots?service_code="+url.QueryEscape(serviceCode), nil, &out)
	switch {
	case err != nil:
		return nil, err
	case status != http.StatusOK:
		return nil, httpjson.Unexpected(service, status)
	}
	return out.toDomain(), nil
}

func (c *Client) Book(ctx context.Context, req domain.AppointmentRequest) (domain.Appointment, error) {
	var out appointmentDTO
	status, err := c.http.Do(ctx, http.MethodPost, c.baseURL+"/api/v1/appointments", newAppointmentRequest(req), &out)
	switch {
	case err != nil:
		return domain.Appointment{}, err
	case status == http.StatusOK, status == http.StatusCreated:
		return out.toDomain(), nil
	case status == http.StatusConflict:
		return domain.Appointment{}, fmt.Errorf("slot %s: %w", req.SlotID, domain.ErrSlotUnavailable)
	case status == http.StatusBadRequest:
		return domain.Appointment{}, fmt.Errorf("%w: slot %s is not available for booking", domain.ErrInvalidInput, req.SlotID)
	default:
		return domain.Appointment{}, httpjson.Unexpected(service, status)
	}
}
