package cli

import (
	"fmt"
	"os"
	"regexp"
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
	fn          func(*rokka.Client, []string)
}

var commands = []Command{
	Command{[]string{"stackoptions", "list"}, "Show default stack options", func(c *rokka.Client, _ []string) {
		res, err := c.GetStackoptions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting stack options: %s", err)
			os.Exit(1)
		}
		PrettyPrintJSON(res)
	}},
	Command{[]string{"organizations", "get", "<name>"}, "Get details of an organization", func(c *rokka.Client, options []string) {
		res, err := c.GetOrganization(options[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting organization: %s", err)
			os.Exit(1)
		}
		PrettyPrintJSON(res)
	}},
}

func GetCommands() []Command {
	return commands
}

func ExecCommand(cl *rokka.Client, userArgs []string) error {
	for _, c := range commands {
		options := []string{}

		for i, arg := range c.Args {
			if len(userArgs) < i+1 {
				break
			}

			isOption, err := regexp.MatchString("^<.*>", arg)
			if err != nil {
				return err
			}

			if isOption {
				options = append(options, userArgs[i])
			} else if arg != userArgs[i] {
				break
			}

			if len(c.Args) == i+1 {
				c.fn(cl, options)
				return nil
			}
		}
	}

	return UnknownCommandError("Unknown command \"" + strings.Join(userArgs, " ") + "\"")
}
