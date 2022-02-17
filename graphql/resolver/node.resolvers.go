package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"giautm.dev/awesome/ent"
	"giautm.dev/awesome/ent/schema/pulid"
)

func (r *queryResolver) Node(ctx context.Context, id pulid.ID) (ent.Noder, error) {
	return r.client.Noder(ctx, id, ent.WithPrefixedULID())
}

func (r *queryResolver) Nodes(ctx context.Context, ids []pulid.ID) ([]ent.Noder, error) {
	return r.client.Noders(ctx, ids, ent.WithPrefixedULID())
}
