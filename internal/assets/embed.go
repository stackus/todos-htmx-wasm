package assets

import (
	"embed"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

//go:embed all:dist
var Assets embed.FS

// Mount mounts the embedded assets to a Chi Router
func Mount(r chi.Router) {
	r.Route("/dist", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasSuffix(r.URL.Path, ".js") {
					w.Header().Add("Service-Worker-Allowed", "/")
				}
				next.ServeHTTP(w, r)
			})
		})
		r.Handle("/*", http.FileServer(http.FS(Assets)))
	})
}
