package pulid

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/vektah/gqlparser/v2/ast"
)

const ulidLength = 26

// Annotation captures the id prefix for a type.
type Annotation struct {
	Prefix string
}

// Name implements the ent Annotation interface.
func (a Annotation) Name() string {
	return "PULID"
}

// MixinWithPrefix creates a Mixin that encodes the provided prefix.
func MixinWithPrefix(prefix string) *Mixin {
	return &Mixin{prefix: prefix}
}

// MixinWithIndex creates a Mixin that encodes the index as base32 for prefix
func MixinWithIndex(idx uint64) *Mixin {
	return MixinWithPrefix(EncodeBase32(idx))
}

// Mixin defines an ent Mixin that captures the PULID prefix for a type.
type Mixin struct {
	mixin.Schema
	prefix string
}

// Annotations returns the annotations for a Mixin instance.
func (m Mixin) Annotations() []schema.Annotation {
	return []schema.Annotation{
		Annotation{Prefix: m.prefix},
		entgql.Directives(
			entgql.NewDirective("pulid", entgql.DirectiveArgument{
				Name:  "prefix",
				Value: m.prefix,
				Kind:  ast.StringValue,
			}),
		),
	}
}

// Fields provides the id field.
func (m Mixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(ID("")).
			MaxLen(len(m.prefix) + 1 + ulidLength).
			DefaultFunc(func() ID {
				return MustNew(m.prefix)
			}),
	}
}
