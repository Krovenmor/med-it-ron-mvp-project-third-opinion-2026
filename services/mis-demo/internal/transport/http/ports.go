package http

import (
	"context"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/b2b"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/mis-demo/internal/demo"
)

type ReportSender interface {
	SendReport(ctx context.Context, r demo.Report) (b2b.Delivery, error)
}
