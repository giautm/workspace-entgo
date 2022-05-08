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

func (r *queryResolver) HelloWorld(ctx context.Context, input model.HelloQueryInput) (*model.HelloQueryResult, error) {
	return &model.HelloQueryResult{
		Message: hello.Hello(input.Name),
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
