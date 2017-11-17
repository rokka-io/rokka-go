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
	Options     []string
	Description string
	fn          func(*rokka.Client, *Log, map[string]string, map[string]string) error
}

func (c *Command) TakesOption(key string) bool {
	for _, v := range c.Options {
		if v == key {
			return true
		}
	}

	return false
}

var Commands = []Command{
	{[]string{"stackoptions", "list"}, nil, "Show default stack options", getStackOptions},
	{[]string{"organizations", "get", "<name>"}, nil, "Get details of an organization", getOrganization},
}

func ExecCommand(cl *rokka.Client, logger *Log, userArgs []string) error {
	hasMatch := false

	for _, c := range Commands {
		commandArgsCount := len(c.Args)

		if len(userArgs) < commandArgsCount {
			continue
		}

		positionalArgs := make(map[string]string)

		for i, arg := range c.Args {
			// check whether this is a positional argument ("<arg>")
			positionalMatch := positionalArgsRegexp.FindStringSubmatch(arg)

			if len(positionalMatch) == 2 {
				positionalArgs[positionalMatch[1]] = userArgs[i]
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
			options := make(map[string]string)

			// parse rest into options ("key=value")
			if len(userArgs) > commandArgsCount {
				for _, option := range userArgs[commandArgsCount-1:] {
					split := strings.Split(option, "=")
					if len(split) == 2 && split[0] != "" && split[1] != "" {
						if !c.TakesOption(split[0]) {
							return fmt.Errorf(`cli: unsupported option "%s" for command "%s"`, split[0], strings.Join(userArgs[:commandArgsCount], " "))
						}

						options[split[0]] = split[1]
					}
				}
			}

			return c.fn(cl, logger, positionalArgs, options)
		}
	}

	return UnknownCommandError(fmt.Sprintf(`cli: Unknown command "%s"`, strings.Join(userArgs, " ")))
}
