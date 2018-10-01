package prisma

import (
	"testing"
)

func TestGetOne(t *testing.T) {
	tests := []struct {
		arg1      interface{}
		arg2      [2]string
		arg3      string
		arg4      []string
		query     string
		numParams int
	}{
		{
			nil, [2]string{"UserWhereUniqueInput!", "User"}, "user", []string{"field1", "field2"},
			"query user {\nuser {\nfield1\nfield2\n}\n}",
			0,
		},
		{
			struct{}{}, [2]string{"UserWhereUniqueInput!", "User"}, "user", []string{"field1", "field2"},
			"query user($where: UserWhereUniqueInput!) {\nuser(where: $where) {\nfield1\nfield2\n}\n}",
			1,
		},
	}

	client := New("")
	for _, tt := range tests {
		exec := client.GetOne(nil, tt.arg1, tt.arg2, tt.arg3, tt.arg4)
		q, params := exec.buildQuery()
		if len(params) != tt.numParams {
			t.Errorf("got %d variables, want %d", len(params), tt.numParams)
		}
		if q != tt.query {
			t.Errorf("got %q, want %q", q, tt.query)
		}
	}
}

func newInt32(v int32) *int32    { return &v }
func newString(v string) *string { return &v }

func TestGetMany(t *testing.T) {
	tests := []struct {
		arg1      *WhereParams
		arg2      [3]string
		arg3      string
		arg4      []string
		query     string
		numParams int
	}{
		{
			&WhereParams{
				Where:   struct{}{},
				OrderBy: newString("ORDERBY"),
				Skip:    newInt32(5),
				After:   newString("AFTER"),
				Before:  newString("BEFORE"),
				First:   newInt32(23),
				Last:    newInt32(42),
			},
			[3]string{"UserWhereInput", "UserOrderByInput", "User"},
			"users",
			[]string{"field1"},
			"query users($where: UserWhereInput, $orderBy: UserOrderByInput, $skip: Int, $after: String, $before: String, $first: Int, $last: Int) {\nusers(where: $where, orderBy: $orderBy, skip: $skip, after: $after, before: $before, first: $first, last: $last) {\nfield1\n}\n}",
			7,
		},
		{
			&WhereParams{
				Where:   struct{}{},
				OrderBy: nil,
				Skip:    newInt32(5),
				After:   nil,
				Before:  newString("BEFORE"),
				First:   newInt32(23),
				Last:    nil,
			},
			[3]string{"UserWhereInput", "UserOrderByInput", "User"},
			"users",
			[]string{"field1"},
			"query users($where: UserWhereInput, $skip: Int, $before: String, $first: Int) {\nusers(where: $where, skip: $skip, before: $before, first: $first) {\nfield1\n}\n}",
			4,
		},
	}

	client := New("")
	for _, tt := range tests {
		exec := client.GetMany(nil, tt.arg1, tt.arg2, tt.arg3, tt.arg4)
		q, params := exec.buildQuery()
		if len(params) != tt.numParams {
			t.Errorf("got %d variables, want %d", len(params), tt.numParams)
		}
		if q != tt.query {
			t.Errorf("got %q, want %q", q, tt.query)
		}
	}
}
