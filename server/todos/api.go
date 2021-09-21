package todos

import (
	"context"
	"errors"

	"github.com/euforic/backend-base/database"
	"github.com/euforic/backend-base/proto"
)

func (s TodosServer) CreateTodo(ctx context.Context, in *proto.CreateTodoReq) (*proto.CreateTodoRes, error) {
	if in == nil {
		return nil, errors.New("missing input")
	}

	res, err := s.db.CreateTodo(ctx, &database.Todo{
		Title:  in.Title,
		Body:   in.Body,
		Author: in.Author,
	})
	if err != nil {
		return nil, err
	}

	out := proto.CreateTodoRes{
		Todo: TodoDBToProto(res),
	}
	return &out, nil
}

func (s TodosServer) GetTodo(ctx context.Context, in *proto.GetTodoReq) (*proto.GetTodoRes, error) {
	var res proto.GetTodoRes

	user, err := s.db.GetTodo(ctx, in.Id)
	res.Todo = TodoDBToProto(user)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s TodosServer) GetTodos(ctx context.Context, in *proto.GetTodosReq) (*proto.GetTodosRes, error) {
	var res proto.GetTodosRes

	for _, id := range in.Ids {
		l, err := s.db.GetTodo(ctx, id)
		if err != nil {
			return nil, err
		}
		res.Todos = append(res.Todos, TodoDBToProto(l))
	}

	return &res, nil
}

func (s TodosServer) UpdateTodo(ctx context.Context, in *proto.UpdateTodoReq) (*proto.UpdateTodoRes, error) {
	var res proto.UpdateTodoRes

	dbres, err := s.db.UpdateTodo(ctx, &database.Todo{
		Id:     in.Id,
		Title:  *in.Title,
		Body:   *in.Body,
		Author: *in.Author,
		IsDone: in.IsDone,
	})
	if err != nil {
		return nil, err
	}

	res.Todo = TodoDBToProto(dbres)
	return &res, nil
}

func (s TodosServer) DeleteTodo(ctx context.Context, in *proto.DeleteTodoReq) (*proto.DeleteTodoRes, error) {
	var res proto.DeleteTodoRes

	if err := s.db.DeleteTodo(ctx, in.Id); err != nil {
		return nil, err
	}
	return &res, nil
}

// ListTodos gRPC list todos
func (s TodosServer) ListTodos(ctx context.Context, in *proto.ListTodosReq) (*proto.ListTodosRes, error) {
	if in.Limit == 0 {
		in.Limit = 20
	}

	var filters database.Filters
	if in.Filters != nil {
		filters = database.Filters{
			Done: in.Filters.Done,
		}
	}

	todos, nextToken, total, err := s.db.ListTodos(ctx, int(in.Offset), int(in.Limit), &filters)
	if err != nil {
		return nil, err
	}

	fds := make([]*proto.Todo, len(todos))
	for i, item := range todos {
		fds[i] = TodoDBToProto(&item)
	}

	return &proto.ListTodosRes{
		Todos:  fds,
		Limit:  int32(len(todos)),
		Offset: int32(nextToken),
		Total:  int32(total),
	}, err
}
