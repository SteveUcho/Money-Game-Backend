package routes

import (
	"context"

	"steveucho.com/packages/backend/gameSystem"
	"steveucho.com/packages/backend/gen/sqlQueries"
)

type App struct {
	DB               *sqlQueries.Queries
	Ctx              context.Context
	GameOrchestrator *gameSystem.GameOrchestrator
}
