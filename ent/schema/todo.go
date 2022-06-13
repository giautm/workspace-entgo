package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"giautm.dev/awesome/ent/schema/pulid"
	"github.com/vektah/gqlparser/v2/ast"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Annotations of the schema.
func (Todo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.Directives(
			entgql.NewDirective("key", entgql.DirectiveArgument{
				Name:  "fields",
				Value: "id",
				Kind:  ast.StringValue,
			}),
		),
		entgql.QueryField(),
		entgql.Mutations(
			entgql.MutationCreate(),
			entgql.MutationUpdate(),
		),
	}
}

// Mixin returns Todo mixed-in schema.
func (Todo) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Time{},
		pulid.MixinWithIndex(845),
	}
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.Text("text").
			NotEmpty().
			Annotations(
				entgql.OrderField("TEXT"),
			),
		field.Enum("status").
			NamedValues(
				"InProgress", "IN_PROGRESS",
				"Completed", "COMPLETED",
			).
			Default("IN_PROGRESS").
			Annotations(
				entgql.OrderField("STATUS"),
			),
		field.Int("priority").
			Default(0).
			Annotations(
				entgql.OrderField("PRIORITY"),
			),
	}
}

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("parent", Todo.Type).
			Unique().
			From("children").
			Annotations(entgql.RelayConnection()),
	}
}
