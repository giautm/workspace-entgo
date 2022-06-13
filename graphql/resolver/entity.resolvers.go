package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"giautm.dev/awesome/ent"
	"giautm.dev/awesome/ent/schema/pulid"
	"giautm.dev/awesome/graphql/generated"
)

func (r *entityResolver) FindTodoByID(ctx context.Context, id pulid.ID) (*ent.Todo, error) {
	todo, err := ent.FromContext(ctx).Todo.Get(ctx, id)
	return todo, ent.MaskNotFound(err)
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
