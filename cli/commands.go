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
	template:    "Id:\t{{.ID}}\nName\t{{.Name}}\nDisplay name:\t{{.DisplayName}}\nBilling email:\t{{.BillingEmail}}\nLimits:\t\n  Space:\t{{.Limit.SpaceInBytes}}\n  Traffic:\t{{.Limit.TrafficInBytes}}\n",
}

var listSourceImagesCommand = Command{
	Args:        []string{"sourceimages", "list", "<org>"},
	QueryParams: []string{"limit", "offset"},
	Description: "List source images",
	fn:          listSourceImages,
	template:    "Name\tHash\tDetails\n{{range .Items}}{{.Name}}\t{{.Hash}}\t{{.MimeType}}, {{.Width}}x{{.Height}}\n{{end}}\nTotal: {{.Total}}\n",
}

var getSourceImageCommand = Command{
	Args:        []string{"sourceimages", "get", "<org>", "<hash>"},
	Description: "Get details of a source image by hash",
	fn:          getSourceImage,
	template:    "Hash:\t{{.Hash}} ({{.ShortHash}})\nName:\t{{.Name}}\nDetails:\t{{.MimeType}}, {{.Width}}x{{.Height}}, {{.Size}}Bytes\nCreated at:\t{{datetime .Created}}\nBinary hash:\t{{.BinaryHash}}\n{{if .UserMetadata}}User metadata:\n{{range $key, $value := .UserMetadata}}  {{$key}}:\t{{$value}}\n{{end}}{{end}}",
}

var getStatsCommand = Command{
	Args:        []string{"stats", "get", "<org>"},
	QueryParams: []string{"from", "to"},
	Description: "Get statistics for an organization",
	fn:          getStats,
	template:    "Date\tDownloaded (Bytes)\tSpace (Bytes)\tFiles\n{{range $_, $e := .}}{{date $e.Date}}\t{{ $e.Downloaded }}\t{{ $e.Space }}\t{{ $e.Files }}\n{{end}}",
}

var listStacksCommand = Command{
	Args:        []string{"stacks", "list", "<org>"},
	Description: "List stacks for an organization",
	fn:          listStacks,
	template:    "Name\tOperations\n{{range .Items}}{{.Name}}\t{{range $i, $e := .StackOperations}}{{if $i}}, {{end}}{{.Name}}{{end}}\n{{end}}",
}

var Commands = []Command{
	getStackOptionsCommand,
	getOrganizationCommand,
	listSourceImagesCommand,
	getSourceImageCommand,
	listStacksCommand,
	getStatsCommand,
}
