package prisma

import (
	"context"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

func (client *Client) decode(exec *Exec, data map[string]interface{}, v interface{}) error {
	var genericData interface{} // This can handle both map[string]interface{} and []interface[]

	// Is unpacking needed
	dataType := reflect.TypeOf(data)
	// XXX this condition is always true, data is statically known to be a map, not an array
	if !IsArray(dataType) {
		unpackedData := data
		for _, instruction := range exec.Stack {
			if IsArray(unpackedData[instruction.Name]) {
				genericData = (unpackedData[instruction.Name]).([]interface{})
				break
			} else {
				unpackedData = (unpackedData[instruction.Name]).(map[string]interface{})
			}
			genericData = unpackedData
		}
	}

	return mapstructure.Decode(genericData, v)
}

func (exec *Exec) Exec(ctx context.Context, v interface{}) error {
	var allArgs []GraphQLArg
	variables := make(map[string]interface{})
	for i := range exec.Stack {
		instruction := &exec.Stack[i]
		for j := range instruction.Args {
			arg := &instruction.Args[j]
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
	query := exec.Client.ProcessInstructions(exec.Stack)
	data, err := exec.Client.GraphQL(ctx, query, variables)
	if err != nil {
		return err
	}

	return exec.Client.decode(exec, data, v)
}

func (exec *Exec) ExecArray(ctx context.Context, v interface{}) error {
	query := exec.Client.ProcessInstructions(exec.Stack)
	variables := make(map[string]interface{})
	for _, instruction := range exec.Stack {
		for _, arg := range instruction.Args {
			variables[arg.Name] = arg.Value
		}
	}
	data, err := exec.Client.GraphQL(ctx, query, variables)
	if err != nil {
		return err
	}
	return exec.Client.decode(exec, data, v)
}
