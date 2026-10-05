package postgres

import (
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/migrations"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/care/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

const SchemaName = "care"

func Schema(q queries.Queries) postgres.Schema {
	return postgres.Schema{Name: SchemaName, Create: q.Schema.Create, Migrations: migrations.Files}
}
