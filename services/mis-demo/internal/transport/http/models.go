package http

import "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/b2b"

type deliveriesResponse struct {
	Delivered []b2b.Delivery `json:"delivered"`
	Error     string         `json:"error,omitempty"`
}
