package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"giautm.dev/awesome/ent"
	"giautm.dev/awesome/ent/schema/pulid"
	"giautm.dev/awesome/ent/todo"
	"giautm.dev/awesome/graphql/generated"
	"giautm.dev/awesome/graphql/model"
	"giautm.dev/awesome/internal/auth"
	"giautm.dev/awesome/internal/hello"
	"giautm.dev/awesome/pkg/logging"
	"go.uber.org/zap"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input ent.CreateTodoInput) (*ent.Todo, error) {
	token := auth.TokenFromContext(ctx)
	if token != nil {
		logger := logging.FromContext(ctx)
		logger.Info("token", zap.String("rawToken", token.Raw))
	}
	return ent.FromContext(ctx).Todo.Create().SetInput(input).Save(ctx)
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, id pulid.ID, input ent.UpdateTodoInput) (*ent.Todo, error) {
	return ent.FromContext(ctx).Todo.UpdateOneID(id).SetInput(input).Save(ctx)
}

func (r *mutationResolver) UpdateTodos(ctx context.Context, ids []pulid.ID, input ent.UpdateTodoInput) ([]*ent.Todo, error) {
	client := ent.FromContext(ctx)
	if err := client.Todo.Update().Where(todo.IDIn(ids...)).SetInput(input).Exec(ctx); err != nil {
		return nil, err
	}
	return client.Todo.Query().Where(todo.IDIn(ids...)).All(ctx)
}

func (r *queryResolver) Todos(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, orderBy *ent.TodoOrder, where *ent.TodoWhereInput) (*ent.TodoConnection, error) {
	return r.client.Todo.Query().
		Paginate(ctx, after, first, before, last,
			ent.WithTodoOrder(orderBy),
			ent.WithTodoFilter(where.Filter),
		)
}

func (r *queryResolver) HelloWorld(ctx context.Context, input model.HelloQueryInput) (*model.HelloQueryResult, error) {
	return &model.HelloQueryResult{
		Message: hello.Hello(input.Name),
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Node(ctx context.Context, id pulid.ID) (ent.Noder, error) {
	return r.client.Noder(ctx, id, ent.WithPrefixedULID())
}
func (r *queryResolver) Nodes(ctx context.Context, ids []pulid.ID) ([]ent.Noder, error) {
	return r.client.Noders(ctx, ids, ent.WithPrefixedULID())
}
