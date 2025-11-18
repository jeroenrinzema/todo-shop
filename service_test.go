package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeroenrinzema/todo-shop/internal/store"
	"github.com/jeroenrinzema/todo-shop/oapi"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGetTODOs(t *testing.T) {
	logger := zaptest.NewLogger(t)
	repository := store.NewRepository()
	service := NewService(logger, repository)

	repository.Set(1, []string{"Buy groceries", "Write tests", "Deploy app"})

	req := httptest.NewRequest(http.MethodGet, "/v1/users/1/todos", nil)
	res := httptest.NewRecorder()

	service.ListTodos(res, req, 1)
	require.Equal(t, http.StatusOK, res.Code)

	var result oapi.TodosList
	err := json.NewDecoder(res.Body).Decode(&result)
	require.NoError(t, err)
	require.Equal(t, 3, result.Total)
	require.Len(t, result.Todos, 3)
	require.EqualValues(t, result.Todos, []string{"Buy groceries", "Write tests", "Deploy app"})
}
