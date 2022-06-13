package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"giautm.dev/awesome/ent"
	"giautm.dev/awesome/ent/schema/pulid"
	"giautm.dev/awesome/graphql/generated"
)

func (r *queryResolver) Node(ctx context.Context, id pulid.ID) (ent.Noder, error) {
	return ent.FromContext(ctx).Noder(ctx, id, ent.WithPrefixedULID())
}

func (r *queryResolver) Nodes(ctx context.Context, ids []pulid.ID) ([]ent.Noder, error) {
	return ent.FromContext(ctx).Noders(ctx, ids, ent.WithPrefixedULID())
}

func (r *queryResolver) Todos(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, orderBy *ent.TodoOrder, where *ent.TodoWhereInput) (*ent.TodoConnection, error) {
	return ent.FromContext(ctx).Todo.Query().
		Paginate(ctx, after, first, before, last,
			ent.WithTodoOrder(orderBy),
			ent.WithTodoFilter(where.Filter),
		)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
