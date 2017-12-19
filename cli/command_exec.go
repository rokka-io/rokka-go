package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/rokka-io/rokka-go/rokka"
)

var positionalArgsRegexp *regexp.Regexp

func init() {
	positionalArgsRegexp = regexp.MustCompile("^<(.*)>$")
}

// UnknownCommandError indicates if the passed arguments to rokka CLI are not known.
type UnknownCommandError string

func (e UnknownCommandError) Error() string {
	return string(e)
}

var funcMap = template.FuncMap{
	"json":     PrettyJSON,
	"datetime": PrettyDateTime,
	"date":     PrettyDate,
}

// ExecCommand searches the available commands, parses arguments accordingly and executes the found command if applicable.
func ExecCommand(cl *rokka.Client, options *CommandOptions, userArgs []string) error {
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
			queryParams := make(map[string]string)

			// parse rest into query params ("key=value")
			if len(userArgs) > commandArgsCount {
				for _, queryParam := range userArgs[commandArgsCount-1:] {
					split := strings.Split(queryParam, "=")
					if len(split) == 2 && split[0] != "" && split[1] != "" {
						if !c.takesQueryParam(split[0]) {
							return fmt.Errorf(`cli: unsupported query parameter "%s" for command "%s"`, split[0], strings.Join(userArgs[:commandArgsCount], " "))
						}

						queryParams[split[0]] = split[1]
					}
				}
			}

			tmpl := c.template
			if len(options.Template) != 0 {
				tmpl = options.Template
			}
			if options.Raw {
				tmpl = rawTemplate
			}

			res, err := c.fn(cl, positionalArgs, queryParams)
			if err != nil {
				return err
			}

			t, err := template.New("").Funcs(funcMap).Parse(tmpl)
			if err != nil {
				return fmt.Errorf("cli: Error parsing response template: %s", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			err = t.Execute(w, res)
			if err != nil {
				return fmt.Errorf("cli: Error formatting response: %s", err)
			}
			w.Flush()

			return nil
		}
	}

	return UnknownCommandError(fmt.Sprintf(`cli: Unknown command "%s"`, strings.Join(userArgs, " ")))
}
