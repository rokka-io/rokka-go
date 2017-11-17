package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"

	"github.com/rokka-io/rokka-go/rokka"
)

func setDefaultValue(data map[string]string, key, val string) {
	if _, ok := data[key]; !ok {
		data[key] = val
	}
}

func extractTemplate(data map[string]string, defaultTemplate string) string {
	setDefaultValue(data, "template", defaultTemplate)
	tmpl := data["template"]
	delete(data, "template")

	return tmpl
}

func getStackOptions(c *rokka.Client, logger *Log, _ map[string]string, _ map[string]string, tmpl *template.Template) error {
	res, err := c.GetStackOptions()
	if err != nil {
		return err
	}
	s, err := PrettyJSON(res)
	if err != nil {
		return fmt.Errorf("cli/getStackOptions: Error pretty printing JSON: %s", err)
	}
	logger.Printf("%s", s)

	return nil
}

func getOrganization(c *rokka.Client, logger *Log, args map[string]string, _ map[string]string, tmpl *template.Template) error {
	res, err := c.GetOrganization(args["name"])
	if err != nil {
		return err
	}
	s, err := PrettyJSON(res)
	if err != nil {
		return fmt.Errorf("cli/getOrganization: Error pretty printing JSON: %s", err)
	}
	logger.Printf("%s", s)

	return nil
}

func listSourceImages(c *rokka.Client, logger *Log, args map[string]string, options map[string]string, tmpl *template.Template) error {
	setDefaultValue(options, "limit", "20")

	res, err := c.ListSourceImages(args["org"], options)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	err = tmpl.Execute(w, res)
	if err != nil {
		return fmt.Errorf("cli/getOrganization: Error formatting response: %s", err)
	}

	w.Flush()

	return nil
}
