package prisma

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type BatchPayloadExec struct {
	client *Client
	stack  []Instruction
}

// BatchPayload docs - generated with types
type BatchPayload struct {
	Count int64 `json:"count"`
}

// Exec docs
func (instance BatchPayloadExec) Exec(ctx context.Context) (BatchPayload, error) {
	var allArgs []GraphQLArg
	variables := make(map[string]interface{})
	for instructionKey := range instance.stack {
		instruction := &instance.stack[instructionKey]
		if instance.client.Debug {
			fmt.Println("Instruction Exec: ", instruction)
		}
		for argKey := range instruction.Args {
			arg := &instruction.Args[argKey]
			if instance.client.Debug {
				fmt.Println("Instruction Arg Exec: ", instruction)
			}
			isUnique := false
			for !isUnique {
				isUnique = true
				for key, existingArg := range allArgs {
					if existingArg.Name == arg.Name {
						isUnique = false
						arg.Name = arg.Name + "_" + strconv.Itoa(key)
						if instance.client.Debug {
							fmt.Println("Resolving Collision Arg Name: ", arg.Name)
						}
						break
					}
				}
			}
			if instance.client.Debug {
				fmt.Println("Arg Name: ", arg.Name)
			}
			allArgs = append(allArgs, *arg)
			variables[arg.Name] = arg.Value
		}
	}
	query := instance.client.ProcessInstructions(instance.stack)
	if instance.client.Debug {
		fmt.Println("Query Exec:", query)
		fmt.Println("Variables Exec:", variables)
	}
	data, err := instance.client.GraphQL(ctx, query, variables)
	if instance.client.Debug {
		fmt.Println("Data Exec:", data)
		fmt.Println("Error Exec:", err)
	}

	var genericData interface{} // This can handle both map[string]interface{} and []interface[]

	// Is unpacking needed
	dataType := reflect.TypeOf(data)
	if !IsArray(dataType) {
		unpackedData := data
		for _, instruction := range instance.stack {
			if instance.client.Debug {
				fmt.Println("Original Unpacked Data Step Exec:", unpackedData)
			}
			if IsArray(unpackedData[instruction.Name]) {
				genericData = (unpackedData[instruction.Name]).([]interface{})
				break
			} else {
				unpackedData = (unpackedData[instruction.Name]).(map[string]interface{})
			}
			if instance.client.Debug {
				fmt.Println("Partially Unpacked Data Step Exec:", unpackedData)
			}
			if instance.client.Debug {
				fmt.Println("Unpacked Data Step Instruction Exec:", instruction.Name)
				fmt.Println("Unpacked Data Step Exec:", unpackedData)
				fmt.Println("Unpacked Data Step Type Exec:", reflect.TypeOf(unpackedData))
			}
			genericData = unpackedData
		}
	}
	if instance.client.Debug {
		fmt.Println("Data Unpacked Exec:", genericData)
	}

	var decodedData BatchPayload
	mapstructure.Decode(genericData, &decodedData)
	if instance.client.Debug {
		fmt.Println("Data Exec Decoded:", decodedData)
	}
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
