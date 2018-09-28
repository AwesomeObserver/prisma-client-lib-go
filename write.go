package prisma

import (
	"context"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type BatchPayloadExec struct {
	client *Client
	stack  []Instruction
}

type BatchPayload struct {
	Count int64 `json:"count"`
}

func (instance BatchPayloadExec) Exec(ctx context.Context) (BatchPayload, error) {
	var allArgs []GraphQLArg
	variables := make(map[string]interface{})
	for instructionKey := range instance.stack {
		instruction := &instance.stack[instructionKey]
		for argKey := range instruction.Args {
			arg := &instruction.Args[argKey]
			isUnique := false
			for !isUnique {
				isUnique = true
				for key, existingArg := range allArgs {
					if existingArg.Name == arg.Name {
						isUnique = false
						arg.Name = arg.Name + "_" + strconv.Itoa(key)
						break
					}
				}
			}
			allArgs = append(allArgs, *arg)
			variables[arg.Name] = arg.Value
		}
	}
	query := instance.client.ProcessInstructions(instance.stack)
	data, err := instance.client.GraphQL(ctx, query, variables)
	if err != nil {
		return BatchPayload{}, err
	}

	var genericData interface{} // This can handle both map[string]interface{} and []interface[]

	// Is unpacking needed
	dataType := reflect.TypeOf(data)
	// XXX this condition is always true, data is statically known to be a map, not an array
	if !isArray(dataType) {
		unpackedData := data
		for _, instruction := range instance.stack {
			if isArray(unpackedData[instruction.Name]) {
				genericData = (unpackedData[instruction.Name]).([]interface{})
				break
			} else {
				unpackedData = (unpackedData[instruction.Name]).(map[string]interface{})
			}
			genericData = unpackedData
		}
	}

	var decodedData BatchPayload
	err = mapstructure.Decode(genericData, &decodedData)
	return decodedData, err
}

type UpdateParams struct {
	Data  interface{}
	Where interface{}
}

func (client *Client) UpdateMany(params UpdateParams, typeNames [2]string, instrName string) *BatchPayloadExec {
	var args []GraphQLArg
	args = append(args, GraphQLArg{
		Name:     "data",
		Key:      "data",
		TypeName: typeNames[0],
		Value:    params.Data,
	})
	if params.Where != nil {
		args = append(args, GraphQLArg{
			Name:     "where",
			Key:      "where",
			TypeName: typeNames[1],
			Value:    params.Where,
		})
	}

	stack := []Instruction{{
		Name: instrName,
		Field: GraphQLField{
			Name:       instrName,
			TypeName:   "BatchPayload",
			TypeFields: []string{"count"},
		},
		Operation: "mutation",
		Args:      args,
	}}

	return &BatchPayloadExec{
		client: client,
		stack:  stack,
	}
}

func (client *Client) Update(params UpdateParams, typeNames [3]string, instrName string, typeFields []string) *Exec {
	var args []GraphQLArg
	args = append(args, GraphQLArg{
		Name:     "data",
		Key:      "data",
		TypeName: typeNames[0],
		Value:    params.Data,
	})
	args = append(args, GraphQLArg{
		Name:     "where",
		Key:      "where",
		TypeName: typeNames[1],
		Value:    params.Where,
	})

	stack := []Instruction{{
		Name: instrName,
		Field: GraphQLField{
			Name:       instrName,
			TypeName:   typeNames[2],
			TypeFields: typeFields,
		},
		Operation: "mutation",
		Args:      args,
	}}

	return &Exec{
		Client: client,
		Stack:  stack,
	}
}

func (client *Client) DeleteMany(params interface{}, typeName string, instrName string) *BatchPayloadExec {
	args := []GraphQLArg{{
		Name:     "where",
		Key:      "where",
		TypeName: typeName,
		Value:    params,
	}}

	stack := []Instruction{{
		Name: instrName,
		Field: GraphQLField{
			Name:       instrName,
			TypeName:   "BatchPayload",
			TypeFields: []string{"count"},
		},
		Operation: "mutation",
		Args:      args,
	}}

	return &BatchPayloadExec{
		client: client,
		stack:  stack,
	}
}

func (client *Client) Delete(params interface{}, typeNames [2]string, instrName string, typeFields []string) *Exec {
	var args []GraphQLArg
	if params != nil {
		args = []GraphQLArg{{
			Name:     "where",
			Key:      "where",
			TypeName: typeNames[0],
			Value:    params,
		}}
	}

	stack := []Instruction{{
		Name: instrName,
		Field: GraphQLField{
			Name:       instrName,
			TypeName:   typeNames[1],
			TypeFields: typeFields,
		},
		Operation: "mutation",
		Args:      args,
	}}

	return &Exec{
		Client: client,
		Stack:  stack,
	}
}

func (client *Client) Create(params interface{}, typeNames [2]string, instrName string, typeFields []string) *Exec {
	var args []GraphQLArg
	if params != nil {
		args = append(args, GraphQLArg{
			Name:     "data",
			Key:      "data",
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
		Operation: "mutation",
		Args:      args,
	}}

	return &Exec{
		Client: client,
		Stack:  stack,
	}
}

type UpsertParams struct {
	Where  interface{}
	Create interface{}
	Update interface{}
}

func (client *Client) Upsert(params *UpsertParams, typeNames [4]string, instrName string, typeFields []string) *Exec {
	var args []GraphQLArg
	if params != nil {
		args = append(args, GraphQLArg{
			Name:     "where",
			Key:      "where",
			TypeName: typeNames[0],
			Value:    params.Where,
		})
		args = append(args, GraphQLArg{
			Name:     "create",
			Key:      "create",
			TypeName: typeNames[1],
			Value:    params.Create,
		})
		args = append(args, GraphQLArg{
			Name:     "update",
			Key:      "update",
			TypeName: typeNames[2],
			Value:    params.Update,
		})
	}

	stack := []Instruction{{
		Name: instrName,
		Field: GraphQLField{
			Name:       instrName,
			TypeName:   typeNames[3],
			TypeFields: typeFields,
		},
		Operation: "mutation",
		Args:      args,
	}}

	return &Exec{
		Client: client,
		Stack:  stack,
	}
}
