package v1

import (
	"encoding/json"
	"github.com/MercyClassic/go_saga/src/app/application/services/command"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/repositories/command"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

func errorWriter(w http.ResponseWriter) {
	http.Error(w, "{\"message\": \"Internal server error\"}", http.StatusInternalServerError)
}

type commandHandler struct {
	userService *commandServices.CommandService
}

func IncludeCommandRouter(r chi.Router, pool client.Client) {
	handler := &commandHandler{
		userService: commandServices.NewCommandService(
			commandRepos.NewCommandRepository(pool),
		),
	}
	r.Get("/commands", handler.getCommands)
	r.Get("/commands/{id}", handler.getCommand)
	r.Post("/commands", handler.createCommand)
}

func (c *commandHandler) getCommands(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := c.userService.GetUsers(ctx)
	if err != nil {
		errorWriter(w)
		return
	}
	err = json.NewEncoder(w).Encode(result)
}

func (c *commandHandler) getCommand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, _ := strconv.Atoi(chi.URLParam(r, "id"))
	result, err := c.userService.GetUser(ctx, userId)
	if err != nil {
		errorWriter(w)
		return
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		errorWriter(w)
		return
	}
}

type CommandCreateRequest struct {
	Description string  `json:"description" validate:"required"`
	UserId      int     `json:"user_id" validate:"required"`
	Amount      float32 `json:"amount" validate:"required"`
}

func (c *commandHandler) createCommand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request CommandCreateRequest
	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := validator.New().Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	result, err := c.userService.CreateUser(
		ctx,
		request.Description,
		request.UserId,
		request.Amount,
	)
	if err != nil {
		errorWriter(w)
		return
	}
	err = json.NewEncoder(w).Encode(result)
}
