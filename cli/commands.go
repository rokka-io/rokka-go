package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/rokka-io/rokka-go/rokka"
)

type UnknownCommandError string

func (e UnknownCommandError) Error() string {
	return string(e)
}

type Command struct {
	Args        []string
	Description string
	fn          func(*rokka.Client)
}

var commands = []Command{
	Command{[]string{"stackoptions", "list"}, "Show default stack options", func(c *rokka.Client) {
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

func ExecCommand(c *rokka.Client, args []string) error {
	command := strings.Join(args, " ")

	for _, v := range commands {
		// TODO: improve this hacky line
		if strings.Join(v.Args, " ") == command {
			v.fn(c)
			return nil
		}
	}

	return UnknownCommandError("Unknown command \"" + command + "\"")
}
