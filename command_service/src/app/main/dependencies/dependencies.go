package dependencies

import (
	"context"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
	"github.com/MercyClassic/go_saga/src/app/presentators/api"
	"github.com/go-chi/chi/v5"
)

func Init(ctx context.Context, r chi.Router, dbUri string) {
	pool, err := client.New(ctx, dbUri)
	if err != nil {
		panic("can't create db pool")
	}

	api.IncludeRouters(r, pool)
}
