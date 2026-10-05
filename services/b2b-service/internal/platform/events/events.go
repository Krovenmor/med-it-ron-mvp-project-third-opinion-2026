package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/jobs"
)

const DeliveryKind = "deliver_event"

type Message struct {
	ID         string
	Topic      string
	Payload    json.RawMessage
	OccurredAt time.Time
}

func (m Message) Decode(dst any) error {
	if err := json.Unmarshal(m.Payload, dst); err != nil {
		return fmt.Errorf("decode %s event: %w", m.Topic, err)
	}
	return nil
}

type Handler func(ctx context.Context, m Message) error

type Subscription struct {
	Topic  string
	Handle Handler
}

type Bus struct {
	handlers map[string][]Handler
}

func NewBus(subscriptions []Subscription) *Bus {
	handlers := make(map[string][]Handler, len(subscriptions))
	for _, s := range subscriptions {
		handlers[s.Topic] = append(handlers[s.Topic], s.Handle)
	}
	return &Bus{handlers: handlers}
}

func (b *Bus) Deliver(ctx context.Context, m Message) error {
	var errs []error
	for _, handle := range b.handlers[m.Topic] {
		if err := handle(ctx, m); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type Clock interface {
	Now() time.Time
}

type Enqueuer interface {
	Enqueue(ctx context.Context, job jobs.Job, now time.Time) error
}

type Publisher struct {
	queue Enqueuer
	clock Clock
}

func NewPublisher(queue Enqueuer, clock Clock) *Publisher {
	return &Publisher{queue: queue, clock: clock}
}

func (p *Publisher) Publish(ctx context.Context, topic, key string, payload any) error {
	now := p.clock.Now()
	job, err := jobs.NewWithPayload(DeliveryKind, key, now, outgoingEnvelope{Topic: topic, Data: payload, OccurredAt: now})
	if err != nil {
		return err
	}
	return p.queue.Enqueue(ctx, job, now)
}

func Delivery(source string, bus *Bus) jobs.Handler {
	return jobs.Handle(DeliveryKind, func(ctx context.Context, job jobs.Job) error {
		var envelope incomingEnvelope
		if err := job.Decode(&envelope); err != nil {
			return err
		}
		return bus.Deliver(ctx, Message{
			ID:         fmt.Sprintf("%s:%d", source, job.ID),
			Topic:      envelope.Topic,
			Payload:    envelope.Data,
			OccurredAt: envelope.OccurredAt,
		})
	})
}
