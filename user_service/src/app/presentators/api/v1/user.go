package v1

import (
	"errors"
	"github.com/MercyClassic/go_saga/src/app/application/services/user"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
	userErrors "github.com/MercyClassic/go_saga/src/app/infrastructure/db/errors"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/repositories/user"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type userHandler struct {
	userService *userServices.UserService
}

func IncludeUserRouter(router *echo.Router, pool client.Client) {
	handler := &userHandler{
		userService: userServices.NewUserService(
			userRepos.NewUserRepository(pool),
		),
	}
	router.Add("GET", "/users", handler.getUsers)
	router.Add("GET", "/users/:id", handler.getUser)
	router.Add("POST", "/users", handler.createUser)
}

func (u *userHandler) getUsers(c echo.Context) error {
	ctx := c.Request().Context()
	result, err := u.userService.GetUsers(ctx)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func (u *userHandler) getUser(c echo.Context) error {
	ctx := c.Request().Context()
	userId, _ := strconv.Atoi(c.Param("id"))
	result, err := u.userService.GetUser(ctx, userId)
	if err != nil {
		if errors.Is(err, userErrors.ErrUserNotFound) {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{"message": userErrors.ErrUserNotFound.Error()},
			)
		}
		return err
	}
	return c.JSON(http.StatusOK, result)
}

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
}

func (u *userHandler) createUser(c echo.Context) error {
	ctx := c.Request().Context()
	var request UserCreateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	result, err := u.userService.CreateUser(
		ctx,
		request.Name,
		request.Username,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, result)
}
