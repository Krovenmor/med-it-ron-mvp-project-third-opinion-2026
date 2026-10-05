package jobs

import (
	"context"
	"time"
)

type Handler interface {
	Kind() string
	Handle(ctx context.Context, job Job) error
}

type HandlerFunc struct {
	kind   string
	handle func(ctx context.Context, job Job) error
}

func Handle(kind string, handle func(ctx context.Context, job Job) error) HandlerFunc {
	return HandlerFunc{kind: kind, handle: handle}
}

func (h HandlerFunc) Kind() string {
	return h.kind
}

func (h HandlerFunc) Handle(ctx context.Context, job Job) error {
	return h.handle(ctx, job)
}

type Clock interface {
	Now() time.Time
}
