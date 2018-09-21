package prisma

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
)

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

func isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func isArray(i interface{}) bool {
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

func (client *Client) ProcessInstructions(stack []Instruction) string {
	query := make(map[string]interface{})
	argsByInstruction := make(map[string][]GraphQLArg)
	var allArgs []GraphQLArg
	firstInstruction := stack[0]
	for i := len(stack) - 1; i >= 0; i-- {
		instruction := stack[i]
		if client.Debug {
			fmt.Println("Instruction: ", instruction)
		}
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
	if client.Debug {
		fmt.Println("Final Query:", query)
		fmt.Println("Final Args By Instruction:", argsByInstruction)
		fmt.Println("Final All Args:", allArgs)
	}
	// TODO: Make this recursive - current depth = 3
	queryTemplateString := `
  {{ $.operation }} {{ $.operationName }} 
  	{{- if eq (len $.allArgs) 0 }} {{ else }} ( {{ end }}
    	{{- range $_, $arg := $.allArgs }}
			\${{ $arg.Name }}: {{ $arg.TypeName }}, 
		{{- end }}
	{{- if eq (len $.allArgs) 0 }} {{ else }} ) {{ end }}
    {
    {{- range $k, $v := $.query }}
    {{- if isArray $v }}
	  {{- $k }}
	  {{- range $argKey, $argValue := $.argsByInstruction }}
	  {{- if eq $argKey $k }}
	  	{{- if eq (len $argValue) 0 }} {{ else }} ( {{ end }}
				{{- range $k, $arg := $argValue}}
					{{ $arg.Key }}: \${{ $arg.Name }},
				{{- end }}
		{{- if eq (len $argValue) 0 }} {{ else }} ) {{ end }}
			{{- end }}
		{{- end }}
	  {
        {{- range $k1, $v1 := $v }}
          {{ $v1 }}
        {{end}}
      }
    {{- else }}
	  {{ $k }} 
	  {{- range $argKey, $argValue := $.argsByInstruction }}
	  	{{- if eq $argKey $k }}
	  		{{- if eq (len $argValue) 0 }} {{ else }} ( {{ end }}
            {{- range $k, $arg := $argValue}}
              {{ $arg.Key }}: \${{ $arg.Name }},
            {{- end }}
			{{- if eq (len $argValue) 0 }} {{ else }} ) {{ end }}
          {{- end }}
        {{- end }}
		{
        {{- range $k, $v := $v }}
        {{- if isArray $v }}
		  {{ $k }} 
		  {{- range $argKey, $argValue := $.argsByInstruction }}
		  {{- if eq $argKey $k }}
			{{- if eq (len $argValue) 0 }} {{ else }} ( {{ end }}
                {{- range $k, $arg := $argValue}}
                  {{ $arg.Key }}: \${{ $arg.Name }},
                {{- end }}
				{{- if eq (len $argValue) 0 }} {{ else }} ) {{ end }} 
              {{- end }}
            {{- end }}
			{ 
            {{- range $k1, $v1 := $v }}
              {{ $v1 }}
            {{end}}
          }
        {{- else }}
		  {{ $k }} 
		  {{- range $argKey, $argValue := $.argsByInstruction }}
		  {{- if eq $argKey $k }}
		  	{{- if eq (len $argValue) 0 }} {{ else }} ( {{ end }}
                {{- range $k, $arg := $argValue}}
                  {{ $arg.Key }}: \${{ $arg.Name }},
                {{- end }}
				{{- if eq (len $argValue) 0 }} {{ else }} ) {{ end }} 
              {{- end }}
            {{- end }}
			{
            {{- range $k, $v := $v }}
              {{- if isArray $v }}
                {{ $k }} { 
                  {{- range $k1, $v1 := $v }}
                    {{ $v1 }}
                  {{end}}
                }
              {{- else }}
				{{ $k }} 
				{{- range $argKey, $argValue := $.argsByInstruction }}
				{{- if eq $argKey $k }}
					{{- if eq (len $argValue) 0 }} {{ else }} ( {{ end }}
                      {{- range $k, $arg := $argValue}}
                        {{ $arg.Key }}: \${{ $arg.Name }},
                      {{- end }}
					  {{- if eq (len $argValue) 0 }} {{ else }} ) {{ end }} 
                    {{- end }}
                  {{- end }}
				  {
                  id
                }
              {{- end }}
              {{- end }}
          }
        {{- end }}
        {{- end }}
      }
    {{- end }}
    {{- end }}
    }
  `
	templateFunctions := template.FuncMap{
		"isArray": isArray,
	}
	queryTemplate, err := template.New("query").Funcs(templateFunctions).Parse(queryTemplateString)
	var queryBytes bytes.Buffer
	var data = make(map[string]interface{})
	data = map[string]interface{}{
		"query":             query,
		"argsByInstruction": argsByInstruction,
		"allArgs":           allArgs,
		"operation":         firstInstruction.Operation,
		"operationName":     firstInstruction.Name,
	}
	queryTemplate.Execute(&queryBytes, data)
	if client.Debug {
		fmt.Println("Query String: ", queryBytes.String())
	}
	if err == nil {
		return queryBytes.String()
	}
	return "Failed to generate query"
}
