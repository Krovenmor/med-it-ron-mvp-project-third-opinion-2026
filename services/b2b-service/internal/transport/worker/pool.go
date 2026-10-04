package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/config"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/domain"
)

type Handler interface {
	Kind() domain.JobKind
	Handle(ctx context.Context, job domain.Job) error
}

type Queue interface {
	Claim(ctx context.Context, now time.Time, lease time.Duration) (domain.Job, bool, error)
	Complete(ctx context.Context, job domain.Job) error
	Reschedule(ctx context.Context, job domain.Job, runAt time.Time, cause string) error
	Fail(ctx context.Context, job domain.Job, cause string) error
}

type Clock interface {
	Now() time.Time
}

type Pool struct {
	cfg      config.Worker
	queue    Queue
	clock    Clock
	log      *zap.Logger
	handlers map[domain.JobKind]Handler
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewPool(cfg config.Worker, queue Queue, clock Clock, log *zap.Logger, handlers []Handler) (*Pool, error) {
	byKind := make(map[domain.JobKind]Handler, len(handlers))
	for _, h := range handlers {
		if _, ok := byKind[h.Kind()]; ok {
			return nil, fmt.Errorf("duplicate handler for job kind %q", h.Kind())
		}
		byKind[h.Kind()] = h
	}
	return &Pool{
		cfg:      cfg,
		queue:    queue,
		clock:    clock,
		log:      log.Named("worker"),
		handlers: byKind,
	}, nil
}

func (p *Pool) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	for range p.cfg.Count {
		p.wg.Go(func() { p.loop(ctx) })
	}
	p.log.Info("worker pool started", zap.Int("workers", p.cfg.Count))
}

func (p *Pool) Stop(ctx context.Context) error {
	p.cancel()
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait for workers: %w", ctx.Err())
	}
}

func (p *Pool) loop(ctx context.Context) {
	for {
		if p.processNext(ctx) {
			continue
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(p.cfg.PollInterval):
		}
	}
}

func (p *Pool) processNext(ctx context.Context) bool {
	job, ok, err := p.queue.Claim(ctx, p.clock.Now(), p.cfg.Lease)
	if err != nil {
		if ctx.Err() == nil {
			p.log.Error("claim job", zap.Error(err))
		}
		return false
	}
	if !ok {
		return false
	}

	jobCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), p.cfg.JobTimeout)
	defer cancel()
	p.settle(jobCtx, job, p.handle(jobCtx, job))
	return true
}

func (p *Pool) handle(ctx context.Context, job domain.Job) (err error) {
	h, ok := p.handlers[job.Kind]
	if !ok {
		return fmt.Errorf("no handler for job kind %q", job.Kind)
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("handler panicked: %v", r)
		}
	}()
	return h.Handle(ctx, job)
}

func (p *Pool) settle(ctx context.Context, job domain.Job, handleErr error) {
	log := p.log.With(
		zap.Int64("job_id", job.ID),
		zap.String("kind", string(job.Kind)),
		zap.Stringer("case_id", job.CaseID),
		zap.Int("attempt", job.Attempts),
	)

	var err error
	switch {
	case handleErr == nil:
		err = p.queue.Complete(ctx, job)
	case job.Attempts >= p.cfg.MaxAttempts:
		log.Error("job failed, attempts exhausted", zap.Error(handleErr))
		err = p.queue.Fail(ctx, job, handleErr.Error())
	default:
		runAt := p.clock.Now().Add(p.backoff(job.Attempts))
		log.Warn("job failed, rescheduled", zap.Error(handleErr), zap.Time("run_at", runAt))
		err = p.queue.Reschedule(ctx, job, runAt, handleErr.Error())
	}

	switch {
	case errors.Is(err, domain.ErrLeaseLost):
		log.Warn("job lease lost before settling")
	case err != nil:
		log.Error("settle job", zap.Error(err))
	}
}

func (p *Pool) backoff(attempt int) time.Duration {
	d := p.cfg.BackoffBase << (attempt - 1)
	if d <= 0 || d > p.cfg.BackoffMax {
		return p.cfg.BackoffMax
	}
	return d
}
