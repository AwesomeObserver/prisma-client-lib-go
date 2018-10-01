package prisma

import "testing"

func TestUpdate(t *testing.T) {
	tests := []struct {
		arg0      UpdateParams
		arg1      [3]string
		arg2      string
		arg3      []string
		query     string
		numParams int
	}{
		{
			UpdateParams{},
			[3]string{"UserUpdateInput!", "UserWhereUniqueInput!", "User"},
			"updateUser",
			[]string{"field1", "field2"},
			"mutation updateUser($data: UserUpdateInput!, $where: UserWhereUniqueInput!) {\nupdateUser(data: $data, where: $where) {\nfield1\nfield2\n}\n}",
			2,
		},
		{
			UpdateParams{
				Data:  struct{}{},
				Where: struct{}{},
			},
			[3]string{"UserUpdateInput!", "UserWhereUniqueInput!", "User"},
			"updateUser",
			[]string{"field1", "field2"},
			"mutation updateUser($data: UserUpdateInput!, $where: UserWhereUniqueInput!) {\nupdateUser(data: $data, where: $where) {\nfield1\nfield2\n}\n}",
			2,
		},
	}

	client := New("")
	for _, tt := range tests {
		exec := client.Update(tt.arg0, tt.arg1, tt.arg2, tt.arg3)
		q, params := exec.buildQuery()
		if len(params) != tt.numParams {
			t.Errorf("got %d variables, want %d", len(params), tt.numParams)
		}
		if q != tt.query {
			t.Errorf("got %q, want %q", q, tt.query)
		}
	}
}
