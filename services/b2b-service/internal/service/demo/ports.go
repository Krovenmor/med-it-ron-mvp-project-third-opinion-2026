package demo

import "context"

type Storage interface {
	Reset(ctx context.Context) error
}

type Clock interface {
	Reset()
}
