package prisma

func (client *Client) GetOne(params interface{}, typeNames [2]string, instrName string, typeFields []string) *Exec {
	var args []GraphQLArg
	if params != nil {
		args = append(args, GraphQLArg{
			Name:     "where",
			Key:      "where",
			TypeName: typeNames[0],
			Value:    params,
		})
	}

	stack := []Instruction{{
		Name: instrName,
		Field: GraphQLField{
			Name:       instrName,
			TypeName:   typeNames[1],
			TypeFields: typeFields,
		},
		Operation: "query",
		Args:      args,
	}}

	return &Exec{
		Client: client,
		Stack:  stack,
	}
}

type WhereParams struct {
	Where   interface{} `json:"where,omitempty"`
	OrderBy *string     `json:"orderBy,omitempty"`
	Skip    *int32      `json:"skip,omitempty"`
	After   *string     `json:"after,omitempty"`
	Before  *string     `json:"before,omitempty"`
	First   *int32      `json:"first,omitempty"`
	Last    *int32      `json:"last,omitempty"`
}

func (client *Client) GetMany(params *WhereParams, typeNames [3]string, instrName string, typeFields []string) *Exec {
	var args []GraphQLArg
	if params != nil {
		if params.Where != nil {
			args = append(args, GraphQLArg{
				Name:     "where",
				Key:      "where",
				TypeName: typeNames[0],
				Value:    params.Where,
			})
		}
		if params.OrderBy != nil {
			args = append(args, GraphQLArg{
				Name:     "orderBy",
				Key:      "orderBy",
				TypeName: typeNames[1],
				Value:    *params.OrderBy,
			})
		}
		if params.Skip != nil {
			args = append(args, GraphQLArg{
				Name:     "skip",
				Key:      "skip",
				TypeName: "Int",
				Value:    *params.Skip,
			})
		}
		if params.After != nil {
			args = append(args, GraphQLArg{
				Name:     "after",
				Key:      "after",
				TypeName: "String",
				Value:    *params.After,
			})
		}
		if params.Before != nil {
			args = append(args, GraphQLArg{
				Name:     "before",
				Key:      "before",
				TypeName: "String",
				Value:    *params.Before,
			})
		}
		if params.First != nil {
			args = append(args, GraphQLArg{
				Name:     "first",
				Key:      "first",
				TypeName: "Int",
				Value:    *params.First,
			})
		}
		if params.Last != nil {
			args = append(args, GraphQLArg{
				Name:     "last",
				Key:      "last",
				TypeName: "Int",
				Value:    *params.Last,
			})
		}
	}

	stack := []Instruction{{
		Name: instrName,
		Field: GraphQLField{
			Name:       instrName,
			TypeName:   typeNames[2],
			TypeFields: typeFields,
		},
		Operation: "query",
		Args:      args,
	}}

	return &Exec{
		Client: client,
		Stack:  stack,
	}
}
