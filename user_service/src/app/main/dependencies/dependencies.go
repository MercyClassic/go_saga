package dependencies

import (
	"context"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
	"github.com/MercyClassic/go_saga/src/app/presentators/api"
	"github.com/labstack/echo/v4"
)

func Init(ctx context.Context, r *echo.Router, dbUri string) {
	pool, err := client.New(ctx, dbUri)
	if err != nil {
		panic("can't create db pool")
	}

	api.IncludeRouters(r, pool)
}
