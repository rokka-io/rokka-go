package cli

import (
	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

func getStackOptions(c *rokka.Client, args []string) (interface{}, error) {
	return c.GetStackOptions()
}

// stackOptionsCmd represents the stackoptions command
var stackOptionsCmd = &cobra.Command{
	Use:                   "stackoptions",
	Short:                 "Inspect stack options",
	Run:                   nil,
	Aliases:               []string{"so"},
	DisableFlagsInUseLine: true,
}

var stackOptionsListCmd = &cobra.Command{
	Use:                   "list",
	Short:                 "Show default stack options",
	Aliases:               []string{"l"},
	DisableFlagsInUseLine: true,
	Run:                   run(getStackOptions, rawTemplate),
}

func init() {
	rootCmd.AddCommand(stackOptionsCmd)

	stackOptionsCmd.AddCommand(stackOptionsListCmd)
}
