package graph

import "graphql/domain/taskstore"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Store taskstore.ITaskStore
}
