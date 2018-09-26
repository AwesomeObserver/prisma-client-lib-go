package prisma

import (
	"context"
	"fmt"
	"reflect"

	"github.com/machinebox/graphql"
)

type Exec struct {
	Client *Client
	Stack  []Instruction
}

type GraphQLField struct {
	Name       string
	TypeName   string
	TypeFields []string
}

type GraphQLArg struct {
	Name     string
	Key      string
	TypeName string
	Value    interface{}
}

type Instruction struct {
	Name      string
	Field     GraphQLField
	Operation string
	Args      []GraphQLArg
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// TODO(dh): get rid of this function if we can
func IsArray(i interface{}) bool {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Array:
		return true
	case reflect.Slice:
		return true
	default:
		return false
	}
}

type Options struct {
	Endpoint string
	Debug    bool
}

func New(options *Options) Client {
	if options == nil {
		return Client{}
	}
	return Client{
		Endpoint: options.Endpoint,
		Debug:    options.Debug,
	}
}

type Client struct {
	Endpoint string
	Debug    bool
}

// GraphQL Send a GraphQL operation request
func (client Client) GraphQL(ctx context.Context, query string, variables map[string]interface{}) (map[string]interface{}, error) {
	// TODO: Add auth support

	req := graphql.NewRequest(query)
	gqlClient := graphql.NewClient(client.Endpoint)

	for key, value := range variables {
		req.Var(key, value)
	}

	// var respData ResponseStruct
	var respData map[string]interface{}
	if err := gqlClient.Run(ctx, req, &respData); err != nil {
		return nil, err
	}
	return respData, nil
}

func (client *Client) ProcessInstructions(stack []Instruction) string {
	query := make(map[string]interface{})
	argsByInstruction := make(map[string][]GraphQLArg)
	var allArgs []GraphQLArg
	firstInstruction := stack[0]

	// XXX why are we walking over the stack backwards? can't we just
	// walk it forwards and construct the AST?
	for i := len(stack) - 1; i >= 0; i-- {
		instruction := stack[i]
		if len(query) == 0 {
			query[instruction.Name] = instruction.Field.TypeFields
			argsByInstruction[instruction.Name] = instruction.Args
			allArgs = append(allArgs, instruction.Args...)
		} else {
			previousInstruction := stack[i+1]
			query[instruction.Name] = map[string]interface{}{
				previousInstruction.Name: query[previousInstruction.Name],
			}
			argsByInstruction[instruction.Name] = instruction.Args
			allArgs = append(allArgs, instruction.Args...)
			delete(query, previousInstruction.Name)
		}
	}

	var opTyp operationType
	switch firstInstruction.Operation {
	case "query":
		opTyp = opQuery
	case "mutation":
		opTyp = opMutation
	case "subscription":
		opTyp = opSubscription
	default:
		// XXX return error
	}
	op := operation{
		typ:  opTyp,
		name: firstInstruction.Name,
	}
	for _, arg := range allArgs {
		op.arguments = append(op.arguments, argument{
			name:  "$" + arg.Name,
			value: arg.TypeName,
		})
	}

	var fn func(root fielder, query map[string]interface{})
	fn = func(root fielder, query map[string]interface{}) {
		// XXX can len(query) ever be larger than 1?
		for k, v := range query {
			q := ObjectField{
				name: k,
			}
			args := argsByInstruction[k]
			for _, arg := range args {
				q.arguments = append(q.arguments, argument{
					name:  arg.Key,
					value: "$" + arg.Name,
				})
			}
			// TODO(dh): redesign the whole instruction processing step,
			// avoid excessive use of interface{} and maps
			switch v := v.(type) {
			case []string:
				for _, f := range v {
					q.fields = append(q.fields, ScalarField{
						name: f,
					})
				}
			case map[string]interface{}:
				fn(&q, v)
			default:
				panic(fmt.Sprintf("unexpected type %T", v))
			}
			root.addField(q)
		}
	}
	fn(&op, query)

	q, err := formatOperation(&op)
	if err != nil {
		// XXX return error
	}
	return q
}
