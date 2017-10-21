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
	fn          func(*rokka.Client, []string, map[string]string)
}

var commands = []Command{
	Command{[]string{"stackoptions", "list"}, "Show default stack options", func(c *rokka.Client, _ []string, _ map[string]string) {
		res, err := c.GetStackoptions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting stack options: %s", err)
			os.Exit(1)
		}
		PrettyPrintJSON(res)
	}},
	Command{[]string{"organizations", "get", "<name>"}, "Get details of an organization", func(c *rokka.Client, args []string, _ map[string]string) {
		res, err := c.GetOrganization(args[0])
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
		positionalArgs := []string{}

		if len(userArgs) < len(c.Args) {
			continue
		}

		for i, arg := range c.Args {
			// check whether this is a positional argument ("<arg>")
			isPositionalArg, err := regexp.MatchString("^<.*>", arg)
			if err != nil {
				return err
			}

			if isPositionalArg {
				positionalArgs = append(positionalArgs, userArgs[i])
			} else if arg != userArgs[i] {
				// user provided argument doesn't match
				break
			}

			// we reached the end
			if len(c.Args) == i+1 {
				options := map[string]string{}

				// parse rest into options ("key=value")
				if len(userArgs) >= i {
					for _, option := range userArgs[i:] {
						split := strings.Split(option, "=")
						if len(split) == 2 && split[0] != "" && split[1] != "" {
							options[split[0]] = split[1]
						}
					}
				}

				c.fn(cl, positionalArgs, options)
				return nil
			}
		}
	}

	return UnknownCommandError("Unknown command \"" + strings.Join(userArgs, " ") + "\"")
}
