package b2b

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

const service = "b2b-service"

type Client struct {
	baseURL      string
	sourceSystem string
	http         httpjson.Client
}

func NewClient(baseURL string, timeout time.Duration, sourceSystem string) *Client {
	return &Client{
		baseURL:      strings.TrimRight(baseURL, "/"),
		sourceSystem: sourceSystem,
		http:         httpjson.New(timeout),
	}
}

func (c *Client) Plan(ctx context.Context, patientID string) (domain.Plan, error) {
	endpoint := c.baseURL + "/api/v1/patients/" + url.PathEscape(patientID) + "/plan?source_system=" + url.QueryEscape(c.sourceSystem)
	var out planResponse
	status, err := c.http.Do(ctx, http.MethodGet, endpoint, nil, &out)
	switch {
	case err != nil:
		return domain.Plan{}, err
	case status == http.StatusNotFound:
		return domain.Plan{}, fmt.Errorf("patient %s: %w", patientID, domain.ErrNotFound)
	case status != http.StatusOK:
		return domain.Plan{}, httpjson.Unexpected(service, status)
	}
	return out.toDomain(), nil
}

func (c *Client) RegisterBooking(ctx context.Context, caseID string, b domain.Booking) error {
	endpoint := c.baseURL + "/api/v1/cases/" + url.PathEscape(caseID) + "/bookings"
	status, err := c.http.Do(ctx, http.MethodPost, endpoint, newBookingRequest(b), nil)
	switch {
	case err != nil:
		return err
	case status == http.StatusOK, status == http.StatusCreated:
		return nil
	case status == http.StatusConflict:
		return fmt.Errorf("%w: b2b-service rejected the booking", domain.ErrInvalidState)
	default:
		return httpjson.Unexpected(service, status)
	}
}

func (c *Client) DeclineRecommendation(ctx context.Context, caseID string, d domain.Decline) error {
	endpoint := c.baseURL + "/api/v1/cases/" + url.PathEscape(caseID) + "/recommendations/" + url.PathEscape(d.RecommendationID) + "/decline"
	status, err := c.http.Do(ctx, http.MethodPost, endpoint, declineRequest{Reason: d.Reason, Comment: d.Comment}, nil)
	switch {
	case err != nil:
		return err
	case status == http.StatusOK:
		return nil
	case status == http.StatusBadRequest:
		return fmt.Errorf("%w: reason must be one of expensive, far, other_clinic, not_needed, other; other requires a comment", domain.ErrInvalidInput)
	case status == http.StatusConflict:
		return fmt.Errorf("%w: the step is already booked or declined", domain.ErrInvalidState)
	case status == http.StatusNotFound:
		return fmt.Errorf("recommendation %s: %w", d.RecommendationID, domain.ErrNotFound)
	}
	return httpjson.Unexpected(service, status)
}

func (c *Client) RequestHelp(ctx context.Context, caseID string) error {
	endpoint := c.baseURL + "/api/v1/cases/" + url.PathEscape(caseID) + "/help-requests"
	status, err := c.http.Do(ctx, http.MethodPost, endpoint, nil, nil)
	switch {
	case err != nil:
		return err
	case status == http.StatusOK, status == http.StatusCreated:
		return nil
	case status == http.StatusNotFound:
		return fmt.Errorf("case %s: %w", caseID, domain.ErrNotFound)
	}
	return httpjson.Unexpected(service, status)
}
