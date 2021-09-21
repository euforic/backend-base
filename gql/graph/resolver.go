package graph

import "github.com/euforic/backend-base/proto"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver struct returns graphQL client
type Resolver struct {
	TodosClient proto.TodosServiceClient
}

// NewResolver ...
// It serves as dependency injection for your app, add any dependencies you require here.
func NewResolver(todos proto.TodosServiceClient) *Resolver {
	return &Resolver{
		TodosClient: todos,
	}
}
