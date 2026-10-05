package demo

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/system/history"
)

type Service struct {
	tx      trm.Manager
	storage Storage
}

func NewService(tx trm.Manager, storage Storage) *Service {
	return &Service{tx: tx, storage: storage}
}

func (s *Service) Reset(ctx context.Context, cases []history.Case) error {
	return s.tx.Do(ctx, func(ctx context.Context) error {
		if err := s.storage.Reset(ctx); err != nil {
			return err
		}
		return s.storage.Import(ctx, cases)
	})
}
