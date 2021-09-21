package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/euforic/backend-base/gql/graph/generated"
	"github.com/euforic/backend-base/gql/graph/model"
	"github.com/euforic/backend-base/pkg/gqltypes"
	"github.com/euforic/backend-base/pkg/ptype"
	"github.com/euforic/backend-base/proto"
)

func (r *dateTimeResolver) Format(ctx context.Context, obj *gqltypes.DateTime, layout *string, timezone *string) (*string, error) {
	var goTime time.Time

	err := obj.UnmarshalGQL(goTime)
	if err != nil {
		return nil, err
	}

	return ptype.StringP(goTime.Format(time.RFC3339)), nil
}

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.CreateTodoInput) (*model.Todo, error) {
	res, err := r.TodosClient.CreateTodo(ctx, &proto.CreateTodoReq{
		Title:  input.Title,
		Body:   input.Body,
		Author: input.Author,
	})
	if err != nil {
		return nil, err
	}

	out := TodoProtoToGQL(res.Todo)
	return out, nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input model.UpdateTodoInput) (*model.Todo, error) {
	res, err := r.TodosClient.UpdateTodo(ctx, &proto.UpdateTodoReq{
		Id:     input.ID,
		Title:  input.Title,
		Body:   input.Body,
		Author: input.Author,
		IsDone: input.IsDone,
	})
	if err != nil {
		return nil, err
	}

	out := TodoProtoToGQL(res.Todo)
	return out, nil
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (bool, error) {
	_, err := r.TodosClient.DeleteTodo(ctx, &proto.DeleteTodoReq{
		Id: id,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Todo(ctx context.Context, id string) (*model.Todo, error) {
	res, err := r.TodosClient.GetTodo(ctx, &proto.GetTodoReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	out := TodoProtoToGQL(res.Todo)
	return out, nil
}

func (r *queryResolver) Todos(ctx context.Context, limit *int, offset *int, filters []*model.TodosFilter) (*model.TodosPayload, error) {
	nodes := []*model.Todo{}

	res, err := r.TodosClient.ListTodos(ctx, &proto.ListTodosReq{
		Offset: int32(ptype.PInt(offset)),
		Limit:  int32(ptype.PInt(limit)),
	})
	if err != nil {
		return nil, err
	}

	for _, l := range res.GetTodos() {
		nodes = append(nodes, TodoProtoToGQL(l))
	}

	out := model.TodosPayload{
		OffsetPageInfo: &model.OffsetPageInfo{
			TotalResults:   int(res.Total),
			Limit:          int(res.Limit),
			Offset:         int(res.Offset),
			NextOffset:     ptype.IntP(int(res.GetOffset() + res.GetLimit())),
			PreviousOffset: ptype.IntP(int(res.GetOffset() - res.GetLimit())),
		},
		Nodes: nodes,
	}
	return &out, nil
}

// DateTime returns generated.DateTimeResolver implementation.
func (r *Resolver) DateTime() generated.DateTimeResolver { return &dateTimeResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type dateTimeResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
