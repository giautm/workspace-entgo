package resolver

import (
	"giautm.dev/awesome/ent"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the GQL resolver root
type Resolver struct {
	client *ent.Client
}

// NewResolver returns a new resolver
func NewResolver(client *ent.Client) *Resolver {
	return &Resolver{
		client: client,
	}
}
