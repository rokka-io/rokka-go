package cli

import (
	"fmt"

	"github.com/rokka-io/rokka-go/rokka"
)

type command struct {
	Args        []string
	QueryParams []string
	Description string
	fn          func(*rokka.Client, map[string]string, map[string]string) (interface{}, error)
	template    string
}

// CommandOptions is a configuration object for the CLI.
type CommandOptions struct {
	Raw      bool
	Template string
}

func (c *command) takesQueryParam(key string) bool {
	for _, v := range c.QueryParams {
		if v == key {
			return true
		}
	}

	return false
}

const rawTemplate = "{{json .}}"

var loginCommand = command{
	Args:        []string{"login"},
	Description: "Store API key in config",
	fn:          login,
	template:    "{{.}}\n",
}

var getStackOptionsCommand = command{
	Args:        []string{"stackoptions", "list"},
	Description: "Show default stack options",
	fn:          getStackOptions,
	template:    rawTemplate,
}

var getOperationsCommand = command{
	Args:        []string{"operations", "list"},
	Description: "Show available operations",
	fn:          getOperations,
	template:    "Name:\tDescription:\tProperties:\n{{range $name, $config := .}}{{$name}}\t{{$config.description}}\t{{join $config.properties \", \"}}\n{{end}}",
}

const organizationTemplate = "Id:\t{{.ID}}\nName\t{{.Name}}\nDisplay name:\t{{.DisplayName}}\nBilling email:\t{{.BillingEmail}}\nLimits:\t\n  Space:\t{{.Limit.SpaceInBytes}}\n  Traffic:\t{{.Limit.TrafficInBytes}}\n"

var getOrganizationCommand = command{
	Args:        []string{"organizations", "get", "<org>"},
	Description: "Get details of an organization",
	fn:          getOrganization,
	template:    organizationTemplate,
}

var createOrganizationCommand = command{
	Args:     []string{"organization", "create", "<org>", "<billingEmail>", "<displayName>"},
	fn:       createOrganization,
	template: organizationTemplate,
}

var createMembershipCommand = command{
	Args:     []string{"membership", "create", "<org>", "<email>", "<role>"},
	fn:       createMembership,
	template: "Membership successfully created.\n",
}

var listSourceImagesCommand = command{
	Args:        []string{"sourceimages", "list", "<org>"},
	QueryParams: []string{"limit", "offset", "hash", "binaryhash", "size", "format", "width", "height", "created"},
	Description: "List source images",
	fn:          listSourceImages,
	template:    "Name\tHash\tDetails\n{{range .Items}}{{.Name}}\t{{.Hash}}\t{{.MimeType}}, {{.Width}}x{{.Height}}\n{{end}}\nTotal: {{.Total}}\n",
}

const sourceImageTemplate = "Hash:\t{{.Hash}} ({{.ShortHash}})\nName:\t{{.Name}}\nDetails:\t{{.MimeType}}, {{.Width}}x{{.Height}}, {{.Size}}Bytes\nCreated at:\t{{datetime .Created}}\nBinary hash:\t{{.BinaryHash}}{{if .UserMetadata}}\nUser metadata:{{range $key, $value := .UserMetadata}}\n  {{$key}}:\t{{$value}}{{end}}{{end}}{{if .DynamicMetadata}}\nDynamic metadata:{{range $key, $value := .DynamicMetadata}}\n  {{$key}}:\t{{$value}}{{end}}{{end}}\n"

var getSourceImageCommand = command{
	Args:        []string{"sourceimages", "get", "<org>", "<hash>"},
	Description: "Get details of a source image by hash",
	fn:          getSourceImage,
	template:    sourceImageTemplate,
}

var createSourceImageCommand = command{
	Args:        []string{"sourceimages", "create", "<org>", "<file>"},
	Description: "Upload image",
	fn:          createSourceImage,
	template:    fmt.Sprintf("{{range .Items}}%s{{end}}", sourceImageTemplate),
}

var getStatsCommand = command{
	Args:        []string{"stats", "get", "<org>"},
	QueryParams: []string{"from", "to"},
	Description: "Get statistics for an organization",
	fn:          getStats,
	template:    "Date\tDownloaded (Bytes)\tSpace (Bytes)\tFiles\n{{range $_, $e := .}}{{date $e.Date}}\t{{ $e.Downloaded }}\t{{ $e.Space }}\t{{ $e.Files }}\n{{end}}",
}

var listStacksCommand = command{
	Args:        []string{"stacks", "list", "<org>"},
	Description: "List stacks for an organization",
	fn:          listStacks,
	template:    "Name\tOperations\n{{range .Items}}{{.Name}}\t{{range $i, $e := .StackOperations}}{{if $i}}, {{end}}{{.Name}}{{end}}\n{{end}}",
}

// Commands list the available and implemented CLI commands.
var Commands = []command{
	loginCommand,
	getStackOptionsCommand,
	getOperationsCommand,
	getOrganizationCommand,
	createOrganizationCommand,
	createMembershipCommand,
	listSourceImagesCommand,
	getSourceImageCommand,
	createSourceImageCommand,
	listStacksCommand,
	getStatsCommand,
}
