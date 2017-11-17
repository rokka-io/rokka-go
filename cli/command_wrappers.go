package cli

import (
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

func getStackOptions(c *rokka.Client, args map[string]string, options map[string]string) (interface{}, error) {
	res, err := c.GetStackOptions()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getOrganization(c *rokka.Client, args map[string]string, options map[string]string) (interface{}, error) {
	res, err := c.GetOrganization(args["name"])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func listSourceImages(c *rokka.Client, args map[string]string, options map[string]string) (interface{}, error) {
	setDefaultValue(options, "limit", "20")

	res, err := c.ListSourceImages(args["org"], options)
	if err != nil {
		return nil, err
	}

	return res, nil
}
