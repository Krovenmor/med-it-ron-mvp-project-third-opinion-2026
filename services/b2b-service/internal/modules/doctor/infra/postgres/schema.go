package postgres

import (
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/infra/postgres/migrations"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/doctor/infra/postgres/queries"
	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/platform/postgres"
)

const SchemaName = "doctor"

func Schema(q queries.Queries) postgres.Schema {
	return postgres.Schema{Name: SchemaName, Create: q.Schema.Create, Migrations: migrations.Files}
}
