package brick

import (
	"testing"
)

func TestSchema(t *testing.T) {
	g := &Graphql{
		schema: "type Query {a(b:B): A} type Mutation {a(b:B): A} type A {ab: ID}",
	}
	g.Schema("type Query {B(a:A): B} type Mutation {a(b:B): A} type B {ab: String}")

	if g.schema != `type Query {
a(b:B):A
B(a:A):B
}
type Mutation {
a(b:B):A
a(b:B):A
}
type A {ab: ID}
type B {ab: String}` {
		t.Error(g.schema)
		t.Error("Schema fail")
	}

	g = &Graphql{
		schema: "type Mutation {a(b:B): A}",
	}
	g.Schema("type Query {B(a:A): B}")

	if g.schema != `type Query {
B(a:A):B
}
type Mutation {
a(b:B):A
}` {
		t.Error(g.schema)
		t.Error("Schema fail")
	}

	g = &Graphql{
		schema: "",
	}
	g.Schema("type Query {B(a:A): B} type Mutation {a(b:B): A}")

	if g.schema != `type Query {
B(a:A):B
}
type Mutation {
a(b:B):A
}` {
		t.Error(g.schema)
		t.Error("Schema fail")
	}

	g = &Graphql{
		schema: "",
	}
	g.Schema("type Mutation {a(b:B): A}")

	if g.schema != `type Mutation {
a(b:B):A
}` {
		t.Error(g.schema)
		t.Error("Schema fail")
	}

	g = &Graphql{
		schema: "type Mutation {a(b:B): A}",
	}
	g.Schema("")

	if g.schema != `type Mutation {
a(b:B):A
}` {
		t.Error(g.schema)
		t.Error("Schema fail")
	}
}
