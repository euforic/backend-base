package buntdb

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/euforic/backend-base/database"
	"github.com/euforic/backend-base/pkg/ptype"
	"github.com/imdario/mergo"
	"github.com/segmentio/ksuid"
	"github.com/tidwall/buntdb"
	"github.com/tidwall/gjson"
)

// CreateTodo new Todo in BuntDB
func (d Database) CreateTodo(ctx context.Context, u *database.Todo) (*database.Todo, error) {
	u.Id = ksuid.New().String()
	u.CreatedAt = time.Now()
	u.IsDone = ptype.BoolP(false)

	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	if err := d.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(todoCollection+":"+u.Id, string(b), nil)
		return err
	}); err != nil {
		return u, err
	}

	return u, nil
}

// GetTodo get Todo by Id
func (d Database) GetTodo(ctx context.Context, Id string) (*database.Todo, error) {
	u := database.Todo{}
	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(todoCollection + ":" + Id)
		if err != nil {
			return err
		}

		deletedAt := gjson.Get(val, "deleted_at").Time()
		if !deletedAt.IsZero() {
			return errors.New("todo not found for id:" + Id)
		}

		if err := json.Unmarshal([]byte(val), &u); err != nil {
			return nil
		}

		return nil
	})

	return &u, err
}

// UpdateTodo updates by Id and accepts a Todo struct for fields to update
func (d Database) UpdateTodo(ctx context.Context, todo *database.Todo) (*database.Todo, error) {
	if todo == nil {
		return nil, errors.New("todo can not be null")
	}
	u := database.Todo{}
	todo.UpdatedAt = time.Now()

	err := d.db.Update(func(tx *buntdb.Tx) error {
		val, err := tx.Get(todoCollection + ":" + todo.Id)
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(val), &u); err != nil {
			return err
		}

		if todo.IsDone != nil {
			u.IsDone = todo.IsDone
		}

		u.UpdatedAt = time.Now()

		if u.DeletedAt != nil && !u.DeletedAt.IsZero() {
			return errors.New("todo not found for id:" + todo.Id)
		}

		if err := mergo.Merge(&u, *todo, mergo.WithOverride); err != nil {
			return err
		}

		udata, err := json.Marshal(u)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(todoCollection+":"+todo.Id, string(udata), nil)
		if err != nil {
			return err
		}

		return nil
	})

	return &u, err
}

// DeleteTodo deletes Todo by Id (as string)
func (d Database) DeleteTodo(ctx context.Context, Id string) error {
	var u database.Todo
	err := d.db.Update(func(tx *buntdb.Tx) error {
		val, err := tx.Get(todoCollection + ":" + Id)
		if err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(val), &u); err != nil {
			return nil
		}

		if u.DeletedAt != nil && !u.DeletedAt.IsZero() {
			return errors.New("todo not found for id:" + Id)
		}

		u.DeletedAt = ptype.PTime(time.Now())

		udata, err := json.Marshal(u)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(todoCollection+":"+u.Id, string(udata), nil)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// ListTodos accepts pagesize and offset params as int and returns a list of todos, new offset and error
func (m *Database) ListTodos(ctx context.Context, offset int, limit int, filters *database.Filters) (todos []database.Todo, newoffset int, total int, err error) {
	var pds []database.Todo
	idx := 0

	err = m.db.View(func(tx *buntdb.Tx) error {
		err := tx.Descend("", func(key, value string) bool {
			if offset > idx {
				idx = idx + 1
				return true
			}

			deletedAt := gjson.Get(value, "deleted_at").Time()
			if !deletedAt.IsZero() {
				return false
			}

			found := true
			if filters != nil {
				if filters.Done != nil {
					if gjson.Get(value, "is_done").Bool() != *filters.Done {
						found = false
					}
				}
			}

			if !found {
				return false
			}

			p := database.Todo{}
			if err = json.Unmarshal([]byte(value), &p); err != nil {
				return false
			}
			pds = append(pds, p)
			idx = idx + 1

			return true
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, 0, 0, err
	}

	newOffset := offset + limit

	// make sure the pageSize doesn't exceed the number of results
	if limit > len(pds) {
		limit = len(pds)
	}
	return pds[0:limit], newOffset, len(pds), nil
}
