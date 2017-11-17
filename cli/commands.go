package cli

import (
	"text/template"

	"github.com/rokka-io/rokka-go/rokka"
)

type Command struct {
	Args        []string
	QueryParams []string
	Description string
	fn          func(*rokka.Client, *Log, map[string]string, map[string]string, *template.Template) error
	template    string
}

type CommandOptions struct {
	Template string
}

func (c *Command) TakesQueryParam(key string) bool {
	if key == "template" {
		return true
	}

	for _, v := range c.QueryParams {
		if v == key {
			return true
		}
	}

	return false
}

var getStackOptionsCommand = Command{
	Args:        []string{"stackoptions", "list"},
	Description: "Show default stack options",
	fn:          getStackOptions,
}
var getOrganizationCommand = Command{
	Args:        []string{"organizations", "get", "<name>"},
	Description: "Get details of an organization",
	fn:          getOrganization,
}

var listStackOptionsCommand = Command{
	Args:        []string{"sourceimages", "list", "<org>"},
	QueryParams: []string{"limit", "offset"},
	Description: "List source images",
	fn:          listSourceImages,
	template:    "Name\tHash\tDetails\n{{range .Items}}{{.Name}}\t{{.Hash}}\t{{.Format}}, {{.Width}}x{{.Height}}\n{{end}}\nTotal: {{.Total}}\n",
}

var Commands = []Command{
	getStackOptionsCommand,
	getOrganizationCommand,
	listStackOptionsCommand,
}
