package todos

import (
	"context"

	"github.com/google/uuid"

	"github.com/stackus/todos-htmx-wasm/internal/domain"
)

type (
	// Service asks "Why do I exist?"; Answer: "To pass the todos."
	Service interface {
		// Add adds a todo to the list
		Add(ctx context.Context, description string) (*domain.Todo, error)
		// Remove removes a todo from the list
		Remove(ctx context.Context, id uuid.UUID) error
		// Update updates a todo in the list
		Update(ctx context.Context, id uuid.UUID, completed bool, description string) (*domain.Todo, error)
		// Search returns a list of todos that match the search string
		Search(ctx context.Context, search string) ([]*domain.Todo, error)
		// Get returns a todo by id
		Get(ctx context.Context, id uuid.UUID) (*domain.Todo, error)
		// Sort sorts the todos by the given ids
		Sort(ctx context.Context, ids []uuid.UUID) ([]*domain.Todo, error)
	}

	service struct {
		todos domain.TodoRepository
	}
)

func NewService(todos domain.TodoRepository) Service {
	return &service{
		todos: todos,
	}
}

func (s service) Add(_ context.Context, description string) (*domain.Todo, error) {
	return s.todos.Add(description)
}

func (s service) Remove(_ context.Context, id uuid.UUID) error {
	return s.todos.Remove(id)
}

func (s service) Update(_ context.Context, id uuid.UUID, completed bool, description string) (*domain.Todo, error) {
	return s.todos.Update(id, completed, description)
}

func (s service) Search(_ context.Context, search string) ([]*domain.Todo, error) {
	return s.todos.Search(search)
}

func (s service) Get(_ context.Context, id uuid.UUID) (*domain.Todo, error) {
	return s.todos.Get(id)
}

func (s service) Sort(_ context.Context, ids []uuid.UUID) ([]*domain.Todo, error) {
	return s.todos.Reorder(ids)
}
