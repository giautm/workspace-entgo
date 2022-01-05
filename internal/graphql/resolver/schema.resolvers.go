package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"giautm.dev/awesome/internal/graphql/generated"
	"giautm.dev/awesome/internal/graphql/model"
	"giautm.dev/awesome/internal/hello"
)

func (r *queryResolver) HelloWorld(ctx context.Context, input model.HelloQueryInput) (*model.HelloQueryResult, error) {
	return &model.HelloQueryResult{
		Message: hello.Hello(input.Name),
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
