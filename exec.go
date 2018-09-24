package prisma

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

func (exec *Exec) Exec(v interface{}) error {
	var allArgs []GraphQLArg
	variables := make(map[string]interface{})
	for i := range exec.Stack {
		instruction := &exec.Stack[i]
		if exec.Client.Debug {
			fmt.Println("Instruction Exec: ", instruction)
		}
		for j := range instruction.Args {
			arg := &instruction.Args[j]
			if exec.Client.Debug {
				fmt.Println("Instruction Arg Exec: ", instruction)
			}
			isUnique := false
			for !isUnique {
				isUnique = true
				for key, existingArg := range allArgs {
					if existingArg.Name == arg.Name {
						isUnique = false
						arg.Name = arg.Name + "_" + strconv.Itoa(key)
						if exec.Client.Debug {
							fmt.Println("Resolving Collision Arg Name: ", arg.Name)
						}
						break
					}
				}
			}
			if exec.Client.Debug {
				fmt.Println("Arg Name: ", arg.Name)
			}
			allArgs = append(allArgs, *arg)
			variables[arg.Name] = arg.Value
		}
	}
	query := exec.Client.ProcessInstructions(exec.Stack)
	if exec.Client.Debug {
		fmt.Println("Query Exec:", query)
		fmt.Println("Variables Exec:", variables)
	}
	data, err := exec.Client.GraphQL(query, variables)
	if exec.Client.Debug {
		fmt.Println("Data Exec:", data)
		fmt.Println("Error Exec:", err)
	}
	if err != nil {
		return err
	}

	var genericData interface{} // This can handle both map[string]interface{} and []interface[]

	// Is unpacking needed
	dataType := reflect.TypeOf(data)
	if !IsArray(dataType) {
		unpackedData := data
		for _, instruction := range exec.Stack {
			if exec.Client.Debug {
				fmt.Println("Original Unpacked Data Step Exec:", unpackedData)
			}
			if IsArray(unpackedData[instruction.Name]) {
				genericData = (unpackedData[instruction.Name]).([]interface{})
				break
			} else {
				unpackedData = (unpackedData[instruction.Name]).(map[string]interface{})
			}
			if exec.Client.Debug {
				fmt.Println("Partially Unpacked Data Step Exec:", unpackedData)
			}
			if exec.Client.Debug {
				fmt.Println("Unpacked Data Step Instruction Exec:", instruction.Name)
				fmt.Println("Unpacked Data Step Exec:", unpackedData)
				fmt.Println("Unpacked Data Step Type Exec:", reflect.TypeOf(unpackedData))
			}
			genericData = unpackedData
		}
	}
	if exec.Client.Debug {
		fmt.Println("Data Unpacked Exec:", genericData)
	}

	err = mapstructure.Decode(genericData, v)
	if exec.Client.Debug {
		fmt.Println("Data Exec Decoded:", v)
	}
	return err
}
