package graph

import (
	"github.com/euforic/backend-base/gql/graph/model"
	"github.com/euforic/backend-base/pkg/gqltypes"
	"github.com/euforic/backend-base/proto"
)

// TodoProtoToGQL ...
func TodoProtoToGQL(in *proto.Todo) *model.Todo {
	if in == nil {
		return &model.Todo{}
	}

	out := model.Todo{
		ID:        in.Id,
		Title:     in.Title,
		Body:      in.Body,
		Author:    in.Author,
		IsDone:    in.IsDone,
		CreatedAt: gqltypes.DateTime(in.CreatedAt.AsTime().String()),
		UpdatedAt: gqltypes.DateTime(in.UpdatedAt.AsTime().String()),
	}

	return &out
}

// TodoGQLToProto ...
func TodoGQLToProto(in *model.Todo) *proto.Todo {
	if in == nil {
		return &proto.Todo{}
	}

	out := proto.Todo{
		Id:     in.ID,
		Title:  in.Title,
		Body:   in.Body,
		Author: in.Author,
		IsDone: in.IsDone,
	}

	return &out
}
