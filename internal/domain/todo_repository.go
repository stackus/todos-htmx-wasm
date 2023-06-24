package domain

import (
	"github.com/google/uuid"
)

type TodoRepository interface {
	Add(description string) (*Todo, error)
	Remove(id uuid.UUID) error
	Update(id uuid.UUID, completed bool, description string) (*Todo, error)
	Search(search string) ([]*Todo, error)
	All() ([]*Todo, error)
	Get(id uuid.UUID) (*Todo, error)
	Reorder(ids []uuid.UUID) ([]*Todo, error)
}
