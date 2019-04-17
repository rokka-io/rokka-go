// This program generates rokka/operations_structs.go
package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/rokka-io/rokka-go/cmd/rokka/cli"
	"github.com/rokka-io/rokka-go/rokka"
)

type operationProperty struct {
	Name string
	Type string
}

type operationProperties []operationProperty

func (o operationProperties) Len() int           { return len(o) }
func (o operationProperties) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o operationProperties) Less(i, j int) bool { return o[i].Name < o[j].Name }

type operation struct {
	Name       string
	Properties operationProperties
	Required   []string
	OneOf      []string
}

type operations []operation

func (o operations) Len() int           { return len(o) }
func (o operations) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o operations) Less(i, j int) bool { return o[i].Name < o[j].Name }

var typeMap = map[string]string{
	"integer": "int",
	"boolean": "bool",
	"string":  "string",
	"number":  "float64",
}

func main() {
	cfg := rokka.Config{}
	generate(&cfg, "operations_objects.go")
}

func generate(cfg *rokka.Config, fileName string) {
	c := rokka.NewClient(cfg)

	res, err := c.GetOperations()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	ops := make(operations, 0)
	for name, value := range res {
		properties := make(operationProperties, 0)
		if propertiesMap, ok := value["properties"].(map[string]interface{}); ok {
			for propName, propValue := range propertiesMap {
				propValueMap := propValue.(map[string]interface{})
				p := operationProperty{
					propName,
					typeMap[propValueMap["type"].(string)],
				}
				properties = append(properties, p)
			}
		}
		sort.Sort(properties)
		var required, oneOf []string
		if list, ok := value["required"]; ok {
			required = cli.ToStringSlice(list.([]interface{}))
		}
		if list, ok := value["oneOf"]; ok {
			oneOf = cli.ToStringSlice(list.([]interface{}))
		}
		o := operation{
			name,
			properties,
			required,
			oneOf,
		}
		ops = append(ops, o)
	}

	sort.Sort(ops)

	var b bytes.Buffer
	err = packageTemplate.Execute(&b, struct {
		Timestamp  time.Time
		Operations operations
	}{
		Timestamp:  time.Now(),
		Operations: ops,
	})
	if err != nil {
		log.Fatal(err)
	}
	src, err := format.Source(b.Bytes())
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		os.Exit(1)
	}
	f.Write(src)
}

var funcMap = template.FuncMap{
	"title":          strings.Title,
	"titleCamelCase": cli.TitleCamelCase,
}

var packageTemplate = template.Must(template.New("").Funcs(funcMap).Parse(`
package rokka

// Code generated by go generate; DO NOT EDIT.
// This file was generated at {{ .Timestamp }}

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Operation is an interface all operation structs implement.
type Operation interface {
	// Name returns the operation's name known by the API.
	Name() string
	// MarshalJSON() implements json.Marshaler 
	MarshalJSON() ([]byte, error)
	// UnmarshalExpressions sets the expressions for an operation
	UnmarshalExpressions(json.RawMessage) error
	// Validate checks if required properties are set.
	// Otherwise it returns false with an error indicating the missing property.
	Validate() (bool, error)
	// toURLPath generates a part of the URL used for dynamic rendering of a stack.
	toURLPath() string
}

type rawStack struct {
	Name    string          ` + "`" + `json:"name"` + "`" + `
	Options json.RawMessage ` + "`" + `json:"options"` + "`" + `
	Expressions json.RawMessage ` + "`" + `json:"expressions"` + "`" + `
}

// Operations is a slice of Operation implementing json.Unmarshaler and json.Marshaler in order to create
// the correct operation types for JSON.
type Operations []Operation

// UnmarshalJSON implements json.Unmarshaler.
func (o *Operations) UnmarshalJSON(data []byte) error {
	ops := make([]rawStack, 0)
	if err := json.Unmarshal(data, &ops); err != nil {
		return err
	}
	for _, v := range ops {
		op, err := NewOperationByName(v.Name)
		if err != nil {
			return err
		}
		*o = append(*o, op.(Operation))
		if err := json.Unmarshal(v.Options, op); err != nil {
			// BUG(mweibel): We continue here when such an error is reached because rokka sometimes (legacy reasons)
			//               has options on an operation which are not of the correct type. Should we write something to stdout? also not nice though..
			continue
		}
		if err := op.UnmarshalExpressions(v.Expressions); err != nil {
			return err
		}
	}
	return nil
}

var errOperationNotImplemented = errors.New("Operation not implemented")

// NewOperationByName creates a struct of the respective type based on the name given.
func NewOperationByName(name string) (Operation, error) {
	switch name {
		{{- range $i, $op := .Operations }}
			case "{{ .Name }}":
				return new({{ title .Name }}Operation), nil
		{{ end }}
	}
	return nil, errOperationNotImplemented
}

{{- range $i, $op := .Operations }}
	// {{ title .Name }}Operation is an auto-generated Operation as specified by the rokka API.
	{{- if (or .OneOf .Required) }}
		// Calling .Validate() will return false if required properties are missing.
	{{- end }}
	//
	// See: https://rokka.io/documentation/references/operations.html
	type {{ title .Name }}Operation struct {
	{{ range .Properties -}}
		{{ titleCamelCase .Name }} *{{ .Type }} ` + "`" + `json:"{{.Name}},omitempty"` + "`" + `
	{{ end }}
		Expressions struct{
			{{ range .Properties -}}
				{{ titleCamelCase .Name }} *string ` + "`" + `json:"{{.Name}},omitempty"` + "`" + `	
			{{ end }}
		} ` + "`" + `json:"expressions,omitempty"` + "`" + `
	}

	// Name implements rokka.Operation.Name
	func (o {{ title .Name }}Operation) Name() string { return "{{ .Name }}" }

	// MarshalJSON() implements rokka.Operation.MarshalJSON
	func (o {{ title .Name }}Operation) MarshalJSON() ([]byte, error) {
		data := make(map[string]interface{})
		data["name"] = o.Name()

		opts := make(map[string]interface{})
		exprs := make(map[string]interface{})

		{{ range .Properties -}}
			if o.{{ titleCamelCase .Name }} != nil {
				opts["{{ .Name }}"] = o.{{ titleCamelCase .Name }}
			}
			if o.Expressions.{{ titleCamelCase .Name }} != nil {
				exprs["{{ .Name }}"] = o.Expressions.{{ titleCamelCase .Name }}
			}
		{{ end }}
		
		if len(opts) > 0 {
			 data["options"] = opts
		}
		if len(exprs) > 0 {
			data["expressions"] = exprs
		}

		return json.Marshal(data)
	}

	// UnmarshalExpressions implements rokka.Operation.UnmarshalExpressions and needs to write back to the struct
	func (o *{{ title .Name }}Operation) UnmarshalExpressions(r json.RawMessage) error {
		return json.Unmarshal(r, &o.Expressions)
	}

	// Validate implements rokka.Operation.Validate.
	func (o {{ title .Name }}Operation) Validate() (bool, error) {
		{{- range .Required }}
			if o.{{ titleCamelCase . }} == nil {
				return false, errors.New("option \"{{ titleCamelCase . }}\" is required")
			}
		{{- end}}
		{{- if .OneOf }}
			valid := false
			{{- range .OneOf }}
				if o.{{ titleCamelCase . }} != nil {
					valid = true
				}
			{{- end}}
			if !valid {
				return false, errors.New("one of \"{{ .OneOf }}\" is required")
			}
		{{- end}}
		return true, nil
	}

	// toURLPath implements rokka.Operation.toURLPath.
	func (o {{ title .Name }}Operation) toURLPath() string {
		options := make([]string, 0)
		{{ range .Properties -}}
			if o.{{ titleCamelCase .Name }} != nil {
			  options = append(options, fmt.Sprintf("%s-%v", "{{ .Name }}", *o.{{ titleCamelCase .Name }}))
			}
		{{ end }}
		if len(options) == 0 {
			return o.Name()
		}
		return fmt.Sprintf("%s-%s", o.Name(), strings.Join(options, "-"))
	}
{{- end }}
`))
