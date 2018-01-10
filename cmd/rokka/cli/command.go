package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

// rawTemplate simply formats a value using prettyJSON.
const rawTemplate = "{{json .}}"

// funcMap contain used helper funcs in the templates.
var funcMap = template.FuncMap{
	"json":     prettyJSON,
	"datetime": prettyDateTime,
	"date":     prettyDate,
	"join":     join,
}

func logErrorAndExit(err error) {
	if err == nil {
		panic("logErrorAndExit called with a nil err")
	}
	logger.Errorf("Error: %s\n", err)
	os.Exit(1)
}

// run returns a func to execute the given func and format the return value with the given template.
func run(fn func(c *rokka.Client, args []string) (interface{}, error), t string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		tmpl := t
		if len(responseTemplate) != 0 {
			tmpl = responseTemplate
		} else if raw {
			tmpl = rawTemplate
		}

		res, err := fn(cl, args)
		if err != nil {
			switch err := err.(type) {
			case *rokka.AnnotatedUnmarshalTypeError:
				logErrorAndExit(fmt.Errorf("%s\n\nRelated JSON:\n----------\n%s\n----------", err, err.Content))
			case rokka.StatusCodeError:
				if err.APIError != nil {
					logErrorAndExit(fmt.Errorf("%s (Status code %d)", err.APIError.Error.Message, err.Code))
				}
				logErrorAndExit(err)
			}
			logErrorAndExit(err)
		}

		t, err := template.New("").Funcs(funcMap).Parse(tmpl)
		if err != nil {
			logErrorAndExit(fmt.Errorf("error parsing response template: %s", err))
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		err = t.Execute(w, res)
		if err != nil {
			logErrorAndExit(fmt.Errorf("cli: Error formatting response: %s", err))
		}
		w.Flush()
	}
}

// prettyJSON marshals JSON with an indent.
func prettyJSON(data interface{}) (string, error) {
	var err error
	var pretty bytes.Buffer

	switch data := data.(type) {
	case []byte:
		err = json.Indent(&pretty, data, "", "  ")
	default:
		enc := json.NewEncoder(&pretty)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		err = enc.Encode(data)
	}

	if err != nil {
		return "", fmt.Errorf("cli/prettyJSON: %s", err)
	}

	return pretty.String(), nil
}

// prettyDateTime formats a time.Time.
func prettyDateTime(t time.Time) string {
	return t.Format("Mon, January _2 2006, 15:04:05Z0700")
}

// prettyDate formats a time.Time.
func prettyDate(t time.Time) string {
	return t.Format("Jan _2 2006")
}

func join(d interface{}, sep string) (string, error) {
	switch d := d.(type) {
	case map[string]interface{}:
		r := make([]string, 0)
		for k := range d {
			r = append(r, k)
		}
		return strings.Join(r, sep), nil
	}
	return "", errors.New("cli/join: unhandled type to join with separator")
}
