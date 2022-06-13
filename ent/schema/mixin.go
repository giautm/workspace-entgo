package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Time composes create/update time mixin.
type Time struct{ mixin.Schema }

// Fields of the time mixin.
func (Time) Fields() []ent.Field {
	return []ent.Field{
		field.Time("create_time").
			Default(time.Now).
			Immutable().
			Annotations(
				entgql.OrderField("CREATE_TIME"),
				entgql.Skip(entgql.SkipMutationCreateInput),
			),
		field.Time("update_time").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(
				entgql.OrderField("UPDATE_TIME"),
				entgql.Skip(entgql.SkipMutationCreateInput),
				entgql.Skip(entgql.SkipMutationUpdateInput),
			),
	}
}
