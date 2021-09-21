package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/euforic/backend-base/database"
	"github.com/imdario/mergo"
	"github.com/segmentio/ksuid"
)

// CreateTodo new Todo in BuntDB
func (d Database) CreateTodo(ctx context.Context, in *database.Todo) (*database.Todo, error) {
	in.Id = ksuid.New().String()
	in.CreatedAt = time.Now()

	s := d.session(ctx)
	if err := s.Create(&in).Error; err != nil {
		return nil, err
	}

	return in, nil
}

// GetTodo get Todo by ID
func (d Database) GetTodo(ctx context.Context, ID string) (*database.Todo, error) {
	todo := database.Todo{}
	s := d.session(ctx)

	if err := s.First(&todo, "id = ?", ID).Error; err != nil {
		s.Logger.Error(ctx, err.Error())
		return nil, err
	}

	return &todo, nil
}

// UpdateTodo updates by ID and accepts a Todo struct for fields to update
func (d Database) UpdateTodo(ctx context.Context, todo *database.Todo) (*database.Todo, error) {
	if todo == nil && todo.Id != "" {
		return nil, errors.New("todo can not be null and must have an ID")
	}

	s := d.session(ctx)

	var updatedTodo database.Todo
	if err := s.First(&updatedTodo, "id = ?", todo.Id).Error; err != nil {
		s.Logger.Error(ctx, err.Error())
		return nil, err
	}
	updatedTodo.UpdatedAt = time.Now()

	if todo.IsDone != nil {
		updatedTodo.IsDone = todo.IsDone
	}

	if err := mergo.Merge(&updatedTodo, todo, mergo.WithOverride); err != nil {
		return nil, err
	}

	if err := s.Updates(updatedTodo).Error; err != nil {
		s.Logger.Error(ctx, err.Error())
		return nil, err
	}

	return &updatedTodo, nil
}

// DeleteTodo deletes Todo by ID (as string)
func (d Database) DeleteTodo(ctx context.Context, ID string) error {
	s := d.session(ctx)

	if err := s.Where("id", ID).Delete(&database.Todo{}).Error; err != nil {
		s.Logger.Error(ctx, err.Error())
		return err
	}
	return nil
}

// ListTodos accepts pagesize and offset params as int and returns a list of todos, new offset and error
func (d *Database) ListTodos(ctx context.Context, offset int, limit int, filters *database.Filters) (todos []database.Todo, newOffset int, total int, err error) {
	s := d.session(ctx).
		Preload("Source")

	if filters != nil {
		if filters.Done != nil {
			if *filters.Done {
				s = s.Or("is_done = ?", *filters.Done)
			}
		}
	}

	err = s.
		Offset(offset).
		Limit(limit).
		Order("todos.created_at desc").
		Find(&todos).
		Error
	if err != nil {
		s.Logger.Error(ctx, err.Error())
		return nil, 0, 0, err
	}

	if len(todos) == int(limit) {
		newOffset = offset + limit
	}

	var c int64
	s.Model(&database.Todo{}).Count(&c)

	return todos, newOffset, int(c), nil
}
