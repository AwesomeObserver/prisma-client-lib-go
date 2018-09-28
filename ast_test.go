package prisma

import (
	"testing"
)

func TestFormatOperation(t *testing.T) {
	tests := []struct {
		op  *operation
		out string
	}{
		{
			&operation{
				typ:  opQuery,
				name: "foo",
			},
			"query foo {\n}",
		},
		{
			&operation{
				typ:       opQuery,
				name:      "foo",
				arguments: []argument{{"$bar", "String!"}},
				fields: []field{
					ScalarField{name: "id"},
					ScalarField{name: "weight", arguments: []argument{{"unit", `"lbs"`}}},
					ObjectField{name: "parent", fields: []field{ScalarField{name: "id"}}},
				},
			},
			"query foo($bar: String!) {\nid\nweight(unit: \"lbs\")\nparent {\nid\n}\n}",
		},
	}

	for _, tt := range tests {
		out := formatOperation(tt.op)
		if out != tt.out {
			t.Errorf("got %q expected %q", out, tt.out)
		}
	}
}

func BenchmarkFormatOperation(b *testing.B) {
	op := &operation{
		typ:       opQuery,
		name:      "foo",
		arguments: []argument{{"$bar", "String!"}},
		fields: []field{
			ScalarField{name: "id"},
			ScalarField{name: "weight", arguments: []argument{{"unit", `"lbs"`}}},
			ObjectField{name: "parent", fields: []field{ScalarField{name: "id"}}},
		},
	}
	for i := 0; i < b.N; i++ {
		formatOperation(op)
	}
}
