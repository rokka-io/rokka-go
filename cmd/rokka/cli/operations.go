package cli

import (
	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

func getOperations(c *rokka.Client, args []string) (interface{}, error) {
	return c.GetOperations()
}

// operationsCmd represents the operations command
var operationsCmd = &cobra.Command{
	Use:                   "operations",
	Short:                 "Inspect operations",
	Run:                   nil,
	Aliases:               []string{"o"},
	DisableFlagsInUseLine: true,
}

var operationsListCmd = &cobra.Command{
	Use:                   "list",
	Short:                 "Show available operations usable for stacks",
	Aliases:               []string{"l"},
	DisableFlagsInUseLine: true,
	Run:                   run(getOperations, "Name:\tDescription:\tProperties:\n{{range $name, $config := .}}{{$name}}\t{{$config.description}}\t{{join $config.properties \", \"}}\n{{end}}"),
}

func init() {
	rootCmd.AddCommand(operationsCmd)

	operationsCmd.AddCommand(operationsListCmd)
}
