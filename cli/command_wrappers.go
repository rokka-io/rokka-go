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

func getSourceImage(c *rokka.Client, args map[string]string, options map[string]string) (interface{}, error) {
	res, err := c.GetSourceImage(args["org"], args["hash"])
	if err != nil {
		return nil, err
	}

	return res, nil
}
