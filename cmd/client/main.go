package main

import (
	"github.com/go-chi/chi/v5"

	wasmhttp "github.com/nlepage/go-wasm-http-server"

	"github.com/stackus/todos-htmx-wasm/cmd/client/htmx"
	"github.com/stackus/todos-htmx-wasm/internal/domain"
	"github.com/stackus/todos-htmx-wasm/internal/features/home"
	"github.com/stackus/todos-htmx-wasm/internal/features/todos"
)

func main() {
	done := make(chan struct{})

	router := chi.NewRouter()
	list := domain.NewTodoApi("http://localhost:3000")

	htmx.Mount(router, htmx.NewHandler(home.NewService(list), todos.NewService(list)))

	println("WASM Client is running")

	wasmhttp.Serve(router)

	// Keep the application running
	<-done
}
