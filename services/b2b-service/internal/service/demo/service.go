package demo

import "context"

type Service struct {
	storage Storage
	clock   Clock
}

func NewService(storage Storage, clock Clock) *Service {
	return &Service{storage: storage, clock: clock}
}

func (s *Service) Reset(ctx context.Context) error {
	if err := s.storage.Reset(ctx); err != nil {
		return err
	}
	s.clock.Reset()
	return nil
}
