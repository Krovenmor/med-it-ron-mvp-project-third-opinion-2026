package clock

import (
	"sync/atomic"
	"time"
)

type Clock struct {
	offset atomic.Int64
}

func New() *Clock {
	return &Clock{}
}

func (c *Clock) Now() time.Time {
	return time.Now().Add(time.Duration(c.offset.Load())).UTC()
}

func (c *Clock) Advance(d time.Duration) time.Time {
	c.offset.Add(int64(d))
	return c.Now()
}
