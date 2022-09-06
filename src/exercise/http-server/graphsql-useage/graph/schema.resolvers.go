package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"graphql/domain/taskstore"
	"graphql/graph/generated"
	"graphql/graph/model"
	"time"
)

// CreateTask is the resolver for the createTask field.
func (r *mutationResolver) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	attachments := make([]*model.Attachment, 0, len(input.Attachments))
	for _, a := range input.Attachments {
		attachments = append(attachments, (*model.Attachment)(a))
	}
	id := r.Store.CreateTask(input.Text, input.Tags, input.Due, []taskstore.Attachment{})
	task := r.Store.GetTask(id)
	return &model.Task{
		ID:          id,
		Text:        task.Text,
		Tags:        task.Tags,
		Due:         task.Due,
		Attachments: attachments,
	}, nil
}

// DeleteTask is the resolver for the deleteTask field.
func (r *mutationResolver) DeleteTask(ctx context.Context, id int) (*bool, error) {
	return nil, r.Store.DeleteTask(id)
}

// DeleteAllTasks is the resolver for the deleteAllTasks field.
func (r *mutationResolver) DeleteAllTasks(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeleteAllTasks - deleteAllTasks"))
}

// GetAllTasks is the resolver for the getAllTasks field.
func (r *queryResolver) GetAllTasks(ctx context.Context) ([]*model.Task, error) {
	tasks := r.Store.GetAllTasks()
	taskDtos := make([]*model.Task, 0, len(tasks))
	for _, task := range tasks {
		taskDtos = append(taskDtos, &model.Task{
			ID:          task.Id,
			Text:        task.Text,
			Tags:        task.Tags,
			Due:         task.Due,
			Attachments: []*model.Attachment{}, // todo
		})
	}
	return taskDtos, nil
}

// GetTask is the resolver for the getTask field.
func (r *queryResolver) GetTask(ctx context.Context, id *int) (*model.Task, error) {
	panic(fmt.Errorf("not implemented: GetTask - getTask"))
}

// GetTasksByTag is the resolver for the getTasksByTag field.
func (r *queryResolver) GetTasksByTag(ctx context.Context, tag string) ([]*model.Task, error) {
	panic(fmt.Errorf("not implemented: GetTasksByTag - getTasksByTag"))
}

// GetTasksByDue is the resolver for the getTasksByDue field.
func (r *queryResolver) GetTasksByDue(ctx context.Context, due time.Time) ([]*model.Task, error) {
	panic(fmt.Errorf("not implemented: GetTasksByDue - getTasksByDue"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
