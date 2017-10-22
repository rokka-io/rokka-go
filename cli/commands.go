package cli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rokka-io/rokka-go/rokka"
)

var positionalArgsRegexp *regexp.Regexp

func init() {
	positionalArgsRegexp = regexp.MustCompile("^<(.*)>$")
}

type UnknownCommandError string

func (e UnknownCommandError) Error() string {
	return string(e)
}

type Command struct {
	Args        []string
	Description string
	fn          func(*rokka.Client, map[string]string, map[string]string) error
}

var Commands = []Command{
	Command{[]string{"stackoptions", "list"}, "Show default stack options", func(c *rokka.Client, _ map[string]string, _ map[string]string) error {
		res, err := c.GetStackoptions()
		if err != nil {
			return err
		}
		fmt.Println(PrettyJSON(res))

		return nil
	}},
	Command{[]string{"organizations", "get", "<name>"}, "Get details of an organization", func(c *rokka.Client, args map[string]string, _ map[string]string) error {
		res, err := c.GetOrganization(args["name"])
		if err != nil {
			return err
		}
		fmt.Println(PrettyJSON(res))

		return nil
	}},
}

func ExecCommand(cl *rokka.Client, userArgs []string) error {
	hasMatch := false

	for _, c := range Commands {
		commandArgsCount := len(c.Args)

		if len(userArgs) < commandArgsCount {
			continue
		}

		positionalArgs := make(map[string]string)

		for i, arg := range c.Args {
			// check whether this is a positional argument ("<arg>")
			positionalMatch := positionalArgsRegexp.FindString(arg)

			if positionalMatch != "" {
				positionalArgs[positionalMatch] = userArgs[i]
			} else if arg != userArgs[i] {
				// user provided argument doesn't match
				break
			}

			// we reached the end
			if commandArgsCount == i+1 {
				hasMatch = true
			}
		}

		if hasMatch {
			options := map[string]string{}

			// parse rest into options ("key=value")
			if len(userArgs) > commandArgsCount {
				for _, option := range userArgs[commandArgsCount-1:] {
					split := strings.Split(option, "=")
					if len(split) == 2 && split[0] != "" && split[1] != "" {
						options[split[0]] = split[1]
					}
				}
			}

			return c.fn(cl, positionalArgs, options)
		}
	}

	return UnknownCommandError(fmt.Sprintf(`Unknown command "%s"`, strings.Join(userArgs, " ")))
}
