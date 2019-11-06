package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

// createStackOverwrite is used for the create stack command, to overwrite it
var createStackOverwrite bool

// stdin is assigned to a variable to make testing easier.
var stdin = os.Stdin

// generatOperationSelection reads the available operations the response of GetOperations
// and transforms that into a list of names and a string to output to the user to select the operations from.
//
// TODO: This would be nicer to handle using the code generator for operations_objects.go because we can only use those
//       anyway for validation etc.
func generateOperationSelection(ops rokka.OperationsResponse) (opNames []string, selectOptions string) {
	for name := range ops {
		opNames = append(opNames, name)
	}
	sort.Strings(opNames)

	prefixes := map[string]bool{
		"q": true,
	}
	for _, name := range opNames {
		prefix := ""
		for i := 1; i < len(name); i++ {
			prefix = string(name[0:i])
			if _, ok := prefixes[prefix]; !ok {
				prefixes[prefix] = true
				selectOptions += fmt.Sprintf("[%s]%s\n", prefix, string(name[len(prefix):]))
				break
			}
		}
	}
	selectOptions += "\n[q]uit\n"

	return
}

// readString reads a string input from stdin and decides if actually something was input.
func readString(r *bufio.Reader) (bool, string, error) {
	c, err := r.ReadString('\n')
	if err != nil {
		return false, "", err
	}
	c = strings.TrimRight(c, "\n")
	if c == "" {
		return true, "", nil
	}
	return false, c, nil
}

func sortedKeysFromMap(m map[string]interface{}) []string {
	l := make([]string, len(m))
	i := 0
	for key := range m {
		l[i] = key
		i++
	}
	sort.Strings(l)
	return l
}

// ToStringSlice is an internally used function to convert a slice of interfaces to a slice of strings
func ToStringSlice(list []interface{}) []string {
	slice := make([]string, len(list))
	for i, v := range list {
		slice[i] = v.(string)
	}
	return slice
}

func shouldRetry(r *bufio.Reader) bool {
	skip, c, err := readString(r)
	if err != nil {
		return false
	}
	if skip || c != "y" {
		return false
	}
	return true
}

func processOptionInput(r *bufio.Reader, props map[string]interface{}, required []string, propName string, s reflect.Value) error {
	fmt.Println()
	options := props[propName].(map[string]interface{})
	fmt.Printf("%s (%s)", propName, options["type"])
	if ListContains(required, propName) {
		fmt.Print(" (required)")
	}
	fmt.Println()
	for name, val := range options {
		fmt.Printf("  %s: %v\n", name, val)
	}
	fmt.Print("Enter value: ")

	skip, c, err := readString(r)
	if err != nil {
		return err
	}
	if skip {
		return nil
	}

	if err := setOperationField(s, propName, c); err != nil {
		fmt.Printf("Input not valid: %s. Retry? [y|n]", err)
		if !shouldRetry(r) {
			return nil
		}
		return processOptionInput(r, props, required, propName, s)
	}
	return nil
}

// processOperationInput generates a rokka.Operation based on the input.
// It requests a new input for every option of that operation. In the end .Validate() is called to check if
// our simple validation check pass. This doesn't guarantee the rokka API will accept the input, though.
func processOperationInput(r *bufio.Reader, name string, props map[string]interface{}, required []string) (rokka.Operation, error) {
	op, err := rokka.NewOperationByName(name)
	if err != nil {
		return nil, err
	}

	ps := reflect.ValueOf(op)
	s := ps.Elem()

	sortedProps := sortedKeysFromMap(props)
	for _, propName := range sortedProps {
		if err := processOptionInput(r, props, required, propName, s); err != nil {
			return nil, err
		}
	}

	if ok, err := op.Validate(); !ok {
		fmt.Printf("Validation failed with error: %s. Retry? [y|n]", err)
		if !shouldRetry(r) {
			return nil, nil
		}
		return processOperationInput(r, name, props, required)
	}
	return op, nil
}

// setOperationField uses reflection to get the type of the field to set, parses the input accordingly,
// and sets the field to that parsed value.
func setOperationField(s reflect.Value, fieldName, val string) error {
	f := s.FieldByName(TitleCamelCase(fieldName))
	switch f.Type().String() {
	case "*int":
		vInt, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		f.Set(reflect.ValueOf(&vInt))
	case "*string":
		f.Set(reflect.ValueOf(&val))
	case "*bool":
		vBool := val == "true"
		f.Set(reflect.ValueOf(&vBool))
	case "*float64":
		vFloat, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		f.Set(reflect.ValueOf(&vFloat))
	}
	return nil
}

// cliCreateStack creates a stack in manual mode.
// Added operations and their options will be added to `req`.
func cliCreateStack(c *rokka.Client, name, org string, req *rokka.CreateStackRequest) error {
	ops, err := c.GetOperations()
	if err != nil {
		return err
	}

	opNames, selectOptions := generateOperationSelection(ops)

	reader := bufio.NewReader(stdin)

	fmt.Printf("Creating stack %s in organization %s\n", name, org)
	for {
		fmt.Printf("Add a new operation by selecting one of the available ops:\n")
		fmt.Println(selectOptions)

		fmt.Print("Enter value: ")
		skip, c, err := readString(reader)
		if err != nil {
			return err
		}
		if skip {
			continue
		}
		if c == "q" {
			break
		}

		for _, name := range opNames {
			if strings.HasPrefix(name, c) {
				opBluePrint := ops[name]
				props := opBluePrint["properties"].(map[string]interface{})
				required := make([]string, 0)
				if opBluePrint["required"] != nil {
					required = ToStringSlice(opBluePrint["required"].([]interface{}))
				}

				op, err := processOperationInput(reader, name, props, required)
				if err != nil {
					return err
				}
				if op != nil {
					req.Operations = append(req.Operations, op)
				}

				break
			}
		}
	}

	return nil
}

// createStack is a slightly more advanced function for creating a stack on the CLI.
// It handles either a manual mode, where the user selects which operation to add and it's parameters, or can also be used in a pipe
// reading JSON data from stdin.
func createStack(c *rokka.Client, args []string) (interface{}, error) {
	org := args[0]
	name := args[1]
	req := rokka.CreateStackRequest{}

	fi, err := stdin.Stat()
	if err != nil {
		return "", err
	}
	if fi.Size() > 0 || fi.Mode()&os.ModeNamedPipe != 0 {
		// data piped in
		if err := json.NewDecoder(stdin).Decode(&req); err != nil {
			return "", err
		}
	} else {
		if err := cliCreateStack(c, name, org, &req); err != nil {
			return "", err
		}
	}

	//v, _ := req.Operations[0].MarshalJSON()
	//fmt.Printf("%s\n", v)

	return c.CreateStack(org, name, req, createStackOverwrite)
}

func listStacks(c *rokka.Client, args []string) (interface{}, error) {
	return c.ListStacks(args[0])
}

func deleteStack(c *rokka.Client, args []string) (interface{}, error) {
	return nil, c.DeleteStack(args[0], args[1])
}

// stacksCmd represents the stacks command
var stacksCmd = &cobra.Command{
	Use:                   "stacks",
	Short:                 "Manage stacks",
	Aliases:               []string{"s"},
	DisableFlagsInUseLine: true,
	Run:                   nil,
}

var stacksListCmd = &cobra.Command{
	Use:                   "list [org]",
	Short:                 "List stacks of an organization",
	Args:                  cobra.ExactArgs(1),
	Aliases:               []string{"l"},
	DisableFlagsInUseLine: true,
	Run:                   run(listStacks, "Name\tOperations\n{{range .Items}}{{.Name}}\t{{range $i, $e := .StackOperations}}{{if $i}}, {{end}}{{.Name}}{{end}}\n{{end}}"),
}

var stacksDeleteCmd = &cobra.Command{
	Use:                   "delete [org] [name]",
	Short:                 "Delete a stack of an organization",
	Args:                  cobra.ExactArgs(2),
	Aliases:               []string{"del"},
	DisableFlagsInUseLine: true,
	Run:                   run(deleteStack, "Stack successfully deleted\n"),
}

var stacksCreateCmd = &cobra.Command{
	Use:   "create [org] [name]",
	Short: "Create or update a stack for an organization",
	Long: `A stack can be created by either passing the JSON data of a new stack in a pipe to this command, or by simply executing the create function.
If the create function is executed without a pipe a manual mode allows to select which operations and their options should be added.`,
	Example: `  # create a stack in manual mode
  rokka stacks create test-organization test-stack

  # create a stack using prefilled JSON data
  echo '{"operations":[{"name":"alpha","options":{"mode":"mask"}}]}' | rokka stacks create test-organization test-stack

  # create a stack using prefilled JSON data from a file
	cat test-stack.json | rokka stacks create test-organization test-stack

	# to update an existing stack, pass the --overwrite flag
	cat test-stack-updated.json | rokka stacks create test-organization test-stack --overwrite`,
	Args:                  cobra.ExactArgs(2),
	Aliases:               []string{"c", "save"},
	DisableFlagsInUseLine: true,
	Run:                   run(createStack, "Stack {{.Name}} created:\n\n{{json .}}\n"),
}

func init() {
	rootCmd.AddCommand(stacksCmd)

	stacksCmd.AddCommand(stacksListCmd)
	stacksCmd.AddCommand(stacksDeleteCmd)
	stacksCmd.AddCommand(stacksCreateCmd)

	stacksCreateCmd.Flags().BoolVar(&createStackOverwrite, "overwrite", false, "Overwrite an existing hash")
}
