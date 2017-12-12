package cli

import (
	"github.com/rokka-io/rokka-go/rokka"
)

func setDefaultValue(data map[string]string, key, val string) {
	if _, ok := data[key]; !ok {
		data[key] = val
	}
}

func getStackOptions(c *rokka.Client, args map[string]string, options map[string]string) (interface{}, error) {
	return c.GetStackOptions()
}

func getOrganization(c *rokka.Client, args map[string]string, options map[string]string) (interface{}, error) {
	return c.GetOrganization(args["name"])
}

func listSourceImages(c *rokka.Client, args map[string]string, options map[string]string) (interface{}, error) {
	setDefaultValue(options, "limit", "20")

	return c.ListSourceImages(args["org"], options)
}

func getSourceImage(c *rokka.Client, args map[string]string, options map[string]string) (interface{}, error) {
	return c.GetSourceImage(args["org"], args["hash"])
}

func listStacks(c *rokka.Client, args map[string]string, options map[string]string) (interface{}, error) {
	return c.ListStacks(args["org"])
}
