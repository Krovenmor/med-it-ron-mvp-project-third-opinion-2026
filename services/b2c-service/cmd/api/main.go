package main

import (
	"go.uber.org/fx"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2c-service/internal/di"
)

func main() {
	fx.New(di.App()).Run()
}
