package main

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	_ "embed"

	"github.com/cloudproud/graceful"
	"github.com/go-chi/chi/v5"
	"github.com/jeroenrinzema/todo-shop/internal/store"
	"github.com/jeroenrinzema/todo-shop/oapi"
	"github.com/jeroenrinzema/todo-shop/pkg/swagger"
	"go.uber.org/zap"
)

//go:embed oapi/resources.yaml
var resources []byte

//go:generate oapi-codegen -package main -o ./oapi/resources_gen.go -generate types,client,chi-server -package oapi oapi/resources.yaml

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	ctx := graceful.NewContext(context.Background())
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	logger.Info("starting TODO shop 🏭")
	store := store.NewRepository()

	router := chi.NewRouter()
	service := NewService(logger, store)

	// NOTE: swagger UI
	root, err := fs.Sub(swagger.FS, "public")
	if err != nil {
		return err
	}

	router.Handle("/public/*", http.FileServer(http.FS(swagger.FS)))
	router.Handle("/openapi.yaml", swagger.HandleOAPI(resources))
	router.Handle("/", http.FileServer(http.FS(root)))

	oapi.HandlerWithOptions(service, oapi.ChiServerOptions{
		BaseRouter: router,
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go server.ListenAndServe()

	logger.Info("server up and running!", zap.String("addr", server.Addr))
	ctx.AwaitKillSignal()
	return nil
}
