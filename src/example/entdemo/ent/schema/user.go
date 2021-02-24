package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	// 添加实体字段
	return []ent.Field{
		field.Int("age").Positive(),
		field.String("name").NotEmpty(),
		field.Bool("sex").Optional(),
		field.String("address"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
	// return []ent.Edge{
	// 	// 一个用户可以有多个 car; 1:N
	// 	edge.To("cars", Car.Type),
	// 	// 一个用户可以在多个组里
	// 	edge.From("groups", Group.Type).
	// 		Ref("users"),
	// }
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "entUsers"},
	}
}
