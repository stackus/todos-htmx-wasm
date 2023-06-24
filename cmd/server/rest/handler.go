package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/segmentio/encoding/json"

	"github.com/stackus/todos-htmx-wasm/internal/features/todos"
	"github.com/stackus/todos-htmx-wasm/internal/log"
)

type (
	Handler interface {
		// All : GET /todos
		All(w http.ResponseWriter, r *http.Request)
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
		todosSvc todos.Service
	}
)

func NewHandler(todosSvc todos.Service) Handler {
	return &handler{todosSvc: todosSvc}
}

func Mount(r chi.Router, h Handler) {
	r.Route("/todos", func(r chi.Router) {
		r.Get("/", h.All)
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

func (h handler) All(w http.ResponseWriter, r *http.Request) {
	todos, err := h.todosSvc.Search(r.Context(), "")
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve todos")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, todos)
}

func (h handler) Sort(w http.ResponseWriter, r *http.Request) {
	type requestType struct {
		Ids []uuid.UUID `json:"ids"`
	}
	var request requestType
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error().Err(err).Msg("failed to decode request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todos, err := h.todosSvc.Sort(r.Context(), request.Ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, todos)
}

func (h handler) Search(w http.ResponseWriter, r *http.Request) {
	var search = r.URL.Query().Get("search")
	todos, err := h.todosSvc.Search(r.Context(), search)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve todos")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, todos)
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	type requestType struct {
		Description string `json:"description"`
	}
	var request requestType
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error().Err(err).Msg("failed to decode request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := h.todosSvc.Add(r.Context(), request.Description)
	if err != nil {
		log.Error().Err(err).Msg("failed to add todo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, todo)
}

func (h handler) Update(w http.ResponseWriter, r *http.Request) {
	type requestType struct {
		Completed   bool   `json:"completed"`
		Description string `json:"description"`
	}
	var request requestType
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error().Err(err).Msg("failed to decode request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id = chi.URLParam(r, "todoId")
	var todoID uuid.UUID
	var err error
	if todoID, err = uuid.Parse(id); err != nil {
		log.Error().Err(err).Msg("failed to parse todoId")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := h.todosSvc.Update(r.Context(), todoID, request.Completed, request.Description)
	if err != nil {
		log.Error().Err(err).Msg("failed to update todo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, todo)
}

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "todoId")
	var todoID uuid.UUID
	var err error
	if todoID, err = uuid.Parse(id); err != nil {
		log.Error().Err(err).Msg("failed to parse todoId")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := h.todosSvc.Get(r.Context(), todoID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get todo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, todo)
}

func (h handler) Delete(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "todoId")
	var todoID uuid.UUID
	var err error
	if todoID, err = uuid.Parse(id); err != nil {
		log.Error().Err(err).Msg("failed to parse todoId")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.todosSvc.Remove(r.Context(), todoID); err != nil {
		log.Error().Err(err).Msg("failed to remove todo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
