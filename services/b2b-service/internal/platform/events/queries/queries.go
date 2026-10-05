package queries

import (
	"embed"
	"fmt"
)

//go:embed *.sql
var files embed.FS

type Queries struct {
	Remember string
}

func Load() (Queries, error) {
	remember, err := files.ReadFile("remember.sql")
	if err != nil {
		return Queries{}, fmt.Errorf("load event queries: %w", err)
	}
	return Queries{Remember: string(remember)}, nil
}
