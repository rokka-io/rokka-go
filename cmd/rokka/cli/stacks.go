package cli

import (
	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

func listStacks(c *rokka.Client, args []string) (interface{}, error) {
	return c.ListStacks(args[0])
}

// stacksCmd represents the stacks command
var stacksCmd = &cobra.Command{
	Use:                   "stacks",
	Short:                 "Manage stacks",
	Aliases:               []string{"s"},
	DisableFlagsInUseLine: true,
	Run: nil,
}

var stacksListCmd = &cobra.Command{
	Use:                   "list [org]",
	Short:                 "List stacks of an organization",
	Args:                  cobra.ExactArgs(1),
	Aliases:               []string{"l"},
	DisableFlagsInUseLine: true,
	Run: run(listStacks, "Name\tOperations\n{{range .Items}}{{.Name}}\t{{range $i, $e := .StackOperations}}{{if $i}}, {{end}}{{.Name}}{{end}}\n{{end}}"),
}

func init() {
	rootCmd.AddCommand(stacksCmd)

	stacksCmd.AddCommand(stacksListCmd)
}
