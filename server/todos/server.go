package todos

import (
	"github.com/euforic/backend-base/database"
	"github.com/euforic/backend-base/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ proto.TodosServiceServer = &TodosServer{}

type TodosServer struct {
	db database.Adapter
}

func New(db database.Adapter) *TodosServer {
	s := TodosServer{db: db}
	return &s
}

func TodoProtoToDB(in *proto.Todo) *database.Todo {
	if in == nil {
		return &database.Todo{}
	}

	out := database.Todo{
		Id:        in.Id,
		Title:     in.Title,
		Body:      in.Body,
		Author:    in.Author,
		IsDone:    &in.IsDone,
		CreatedAt: in.CreatedAt.AsTime(),
		UpdatedAt: in.UpdatedAt.AsTime(),
	}

	return &out
}

func TodoDBToProto(in *database.Todo) *proto.Todo {
	if in == nil {
		return &proto.Todo{}
	}

	out := proto.Todo{
		Id:        in.Id,
		Title:     in.Title,
		Body:      in.Body,
		Author:    in.Author,
		IsDone:    *in.IsDone,
		CreatedAt: timestamppb.New(in.CreatedAt),
		UpdatedAt: timestamppb.New(in.UpdatedAt),
	}

	if in.DeletedAt != nil {
		out.DeletedAt = timestamppb.New(*in.DeletedAt)
	}

	return &out
}
