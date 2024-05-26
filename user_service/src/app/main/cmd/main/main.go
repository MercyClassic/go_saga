package main

import (
	"context"
	"github.com/MercyClassic/go_saga/src/app/main/dependencies"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	ctx := context.Background()
	dependencies.Init(
		ctx,
		e.Router(),
		os.Getenv("db_uri"),
	)
	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
