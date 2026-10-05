package http

import "time"

type advanceClockRequest struct {
	By string `json:"by"`
}

type clockResponse struct {
	Now time.Time `json:"now"`
}
