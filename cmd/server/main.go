package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/stackus/todos-htmx-wasm/cmd/server/rest"
	"github.com/stackus/todos-htmx-wasm/internal/assets"
	"github.com/stackus/todos-htmx-wasm/internal/domain"
	"github.com/stackus/todos-htmx-wasm/internal/features/todos"
	"github.com/stackus/todos-htmx-wasm/internal/log"
	"github.com/stackus/todos-htmx-wasm/internal/templates/pages"
)

func main() {
	var port = ":3000"

	flag.StringVar(&port, "port", port, "port to listen on")
	flag.Parse()

	router := chi.NewRouter()
	router.Use(
		log.WebLogger(log.DefaultLogger),
		middleware.Recoverer,
		middleware.Compress(5),
	)

	list := domain.NewTodos()
	list.Add("Bake a cake")
	list.Add("Feed the cat")
	list.Add("Take out the trash")

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := pages.LoadingPage().Render(r.Context(), w); err != nil {
			log.Error().Err(err).Msg("failed to render loading page")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	rest.Mount(router, rest.NewHandler(todos.NewService(list)))
	assets.Mount(router)

	server := &http.Server{
		Addr:    port,
		Handler: http.TimeoutHandler(router, 30*time.Second, "request timed out"),
	}

	// Display the localhost address and port
	fmt.Printf("Listening on http://localhost%s\n", port)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
