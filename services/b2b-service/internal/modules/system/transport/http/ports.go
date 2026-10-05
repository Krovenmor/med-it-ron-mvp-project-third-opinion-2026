package http

import (
	"context"
	"time"
)

type Clock interface {
	Now() time.Time
	Advance(d time.Duration) time.Time
}

type Demo interface {
	Reset(ctx context.Context) error
}
