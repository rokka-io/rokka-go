package cli

import (
	"github.com/rokka-io/rokka-go/rokka"
)

type Command struct {
	Args        []string
	QueryParams []string
	Description string
	fn          func(*rokka.Client, map[string]string, map[string]string) (interface{}, error)
	template    string
}

type CommandOptions struct {
	Raw      bool
	Template string
}

func (c *Command) TakesQueryParam(key string) bool {
	for _, v := range c.QueryParams {
		if v == key {
			return true
		}
	}

	return false
}

const rawTemplate = "{{json .}}"

var getStackOptionsCommand = Command{
	Args:        []string{"stackoptions", "list"},
	Description: "Show default stack options",
	fn:          getStackOptions,
	template:    rawTemplate,
}
var getOrganizationCommand = Command{
	Args:        []string{"organizations", "get", "<name>"},
	Description: "Get details of an organization",
	fn:          getOrganization,
	template:    "Id:\t{{.ID}}\nName\t{{.Name}}\nDisplay name:\t{{.DisplayName}}\nBilling email:\t{{.BillingEmail}}\nLimits:\t\n  Space:\t{{.Limit.SpaceInBytes}}\n  Traffic:\t{{.Limit.TrafficInBytes}}",
}

var listStackOptionsCommand = Command{
	Args:        []string{"sourceimages", "list", "<org>"},
	QueryParams: []string{"limit", "offset"},
	Description: "List source images",
	fn:          listSourceImages,
	template:    "Name\tHash\tDetails\n{{range .Items}}{{.Name}}\t{{.Hash}}\t{{.Format}}, {{.Width}}x{{.Height}}\n{{end}}\nTotal: {{.Total}}",
}

var Commands = []Command{
	getStackOptionsCommand,
	getOrganizationCommand,
	listStackOptionsCommand,
}
