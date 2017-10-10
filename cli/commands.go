package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/rokka-io/rokka-go/client"
)

type Command struct {
	args        []string
	description string
	fn          func(*client.Client)
}

var commands = []Command{
	Command{[]string{"stackoptions", "list"}, "Show default stack options", func(c *client.Client) {
		res, err := c.GetStackoptions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting stack options: %s", err)
			os.Exit(1)
		}
		PrettyPrintJSON(res)
	}},
}

func GetCommands() []Command {
	return commands
}

func ExecCommand(c *client.Client, args []string) {
	for _, v := range commands {
		// TODO: improve this hacky line
		if strings.Join(v.args, " ") == strings.Join(args, " ") {
			v.fn(c)
			return
		}
	}
}
