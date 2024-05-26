package api

import (
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
	"github.com/MercyClassic/go_saga/src/app/presentators/api/v1"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}

func IncludeRouters(r *echo.Router, pool client.Client) {
	r.Add("GET", "/ping", ping)
	v1.IncludeUserRouter(r, pool)
}
