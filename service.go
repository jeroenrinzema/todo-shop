package main

import (
	"net/http"

	"github.com/jeroenrinzema/todo-shop/internal/json"
	"github.com/jeroenrinzema/todo-shop/internal/store"
	"github.com/jeroenrinzema/todo-shop/oapi"
	"go.uber.org/zap"
)

func NewService(logger *zap.Logger, conn *store.Repository) *Service {
	return &Service{
		logger: logger,
		store:  conn,
	}
}

type Service struct {
	logger *zap.Logger
	store  *store.Repository
}

func (svc *Service) ListTodos(w http.ResponseWriter, r *http.Request, userId int) {
	logger := svc.logger.With(zap.Int("user_id", userId))
	logger.Info("listing todos")

	todos := svc.store.Get(userId)
	res := oapi.TodosList{
		Todos: todos,
		Total: len(todos),
	}

	logger.Info("listed todos")
	json.Write(w, http.StatusOK, res)
}

func (svc *Service) CreateTodo(w http.ResponseWriter, r *http.Request, userId int) {
	logger := svc.logger.With(zap.Int("user_id", userId))
	logger.Info("creating todo")

	var req oapi.CreateTodoJSONRequestBody
	if err := json.Decode(r.Body, &req); err != nil {
		logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	svc.store.Append(userId, req.Description)

	logger.Info("created todo")
	json.Write(w, http.StatusCreated, req.Description)
}

func (svc *Service) DeleteTodo(w http.ResponseWriter, r *http.Request, userId int, todoId int) {
	// TODO: to be implemented!
}
