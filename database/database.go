package database

import "context"

type Adapter interface {
	Connect() error
	Close() error

	CreateTodo(ctx context.Context, in *Todo) (*Todo, error)
	GetTodo(ctx context.Context, id string) (*Todo, error)
	UpdateTodo(ctx context.Context, in *Todo) (*Todo, error)
	DeleteTodo(ctx context.Context, id string) error
	ListTodos(ctx context.Context, offset int, limit int, filters *Filters) ([]Todo, int, int, error)
}

type Filters struct {
	Done *bool `json:"done"`
}
