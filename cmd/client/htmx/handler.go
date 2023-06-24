package htmx

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/stackus/todos-htmx-wasm/internal/features/home"
	"github.com/stackus/todos-htmx-wasm/internal/features/todos"
	"github.com/stackus/todos-htmx-wasm/internal/templates/pages"
	"github.com/stackus/todos-htmx-wasm/internal/templates/partials"
)

type (
	Handler interface {
		// Home : GET /
		Home(w http.ResponseWriter, r *http.Request)
		// Search : GET /todos
		Search(w http.ResponseWriter, r *http.Request)
		// Create : POST /todos
		Create(w http.ResponseWriter, r *http.Request)
		// Update : PATCH /todos/{todoId}
		// Update : POST /todos/{todoId}/edit
		Update(w http.ResponseWriter, r *http.Request)
		// Get : GET /todos/{todoId}
		Get(w http.ResponseWriter, r *http.Request)
		// Delete : DELETE /todos/{todoId}
		// Delete : POST /todos/{todoId}/delete
		Delete(w http.ResponseWriter, r *http.Request)
		// Sort : POST /todos/sort
		Sort(w http.ResponseWriter, r *http.Request)
	}

	handler struct {
		homeSvc  home.Service
		todosSvc todos.Service
	}
)

func NewHandler(homeSvc home.Service, todosSvc todos.Service) Handler {
	return &handler{
		homeSvc:  homeSvc,
		todosSvc: todosSvc,
	}
}

func Mount(r chi.Router, h Handler) {
	r.Get("/", h.Home)
	r.Route("/todos", func(r chi.Router) {
		r.Get("/", h.Search)
		r.Post("/", h.Create)
		r.Route("/{todoId}", func(r chi.Router) {
			r.Patch("/", h.Update)
			r.Post("/edit", h.Update)
			r.Get("/", h.Get)
			r.Delete("/", h.Delete)
			r.Post("/delete", h.Delete)
		})
		r.Post("/sort", h.Sort)
	})
}

func (h handler) Home(w http.ResponseWriter, r *http.Request) {
	list, err := h.homeSvc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := pages.HomePage(list).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) Sort(w http.ResponseWriter, r *http.Request) {
	var todoIDs []uuid.UUID
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, id := range r.Form["id"] {
		var todoID uuid.UUID
		var err error
		if todoID, err = uuid.Parse(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		todoIDs = append(todoIDs, todoID)
	}
	if _, err := h.todosSvc.Sort(r.Context(), todoIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h handler) Search(w http.ResponseWriter, r *http.Request) {
	var search = r.URL.Query().Get("search")
	list, err := h.todosSvc.Search(r.Context(), search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		err = partials.RenderTodos(list).Render(r.Context(), w)
	default:
		err = pages.TodosPage(list, search).Render(r.Context(), w)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var description = r.Form.Get("description")

	todo, err := h.todosSvc.Add(r.Context(), description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		err = partials.RenderTodo(todo).Render(r.Context(), w)
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) Update(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "todoId")
	var todoID uuid.UUID
	var err error
	if todoID, err = uuid.Parse(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var completed = r.Form.Get("completed") == "true"
	var description = r.Form.Get("description")

	todo, err := h.todosSvc.Update(r.Context(), todoID, completed, description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		err = partials.RenderTodo(todo).Render(r.Context(), w)
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "todoId")
	var todoID uuid.UUID
	var err error
	if todoID, err = uuid.Parse(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo, err := h.todosSvc.Get(r.Context(), todoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		err = partials.EditTodoForm(todo).Render(r.Context(), w)
	default:
		err = pages.TodoPage(todo).Render(r.Context(), w)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) Delete(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "todoId")
	var todoID uuid.UUID
	var err error
	if todoID, err = uuid.Parse(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.todosSvc.Remove(r.Context(), todoID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch isHTMX(r) {
	case true:
		_, err = w.Write([]byte(""))
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func isHTMX(r *http.Request) bool {
	// Check for "HX-Request" header
	if r.Header.Get("HX-Request") != "" {
		return true
	}

	return false
}
