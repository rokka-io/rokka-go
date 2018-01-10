package cli

import (
	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

func getOrganization(c *rokka.Client, args []string) (interface{}, error) {
	return c.GetOrganization(args[0])
}

func createOrganization(c *rokka.Client, args []string) (interface{}, error) {
	return c.CreateOrganization(args[0], args[1], args[2])
}

// organizationCmd represents the stackoptions command
var organizationCmd = &cobra.Command{
	Use:                   "organization",
	Short:                 "Manage organizations",
	Run:                   nil,
	Aliases:               []string{"org"},
	DisableFlagsInUseLine: true,
}

const organizationTemplate = "Id:\t{{.ID}}\nName\t{{.Name}}\nDisplay name:\t{{.DisplayName}}\nBilling email:\t{{.BillingEmail}}\nLimits:\t\n  Space:\t{{.Limit.SpaceInBytes}}\n  Traffic:\t{{.Limit.TrafficInBytes}}\n"

var organizationGetCmd = &cobra.Command{
	Use:                   "get [org]",
	Short:                 "Get details of an organization",
	Args:                  cobra.ExactArgs(1),
	Aliases:               []string{"g"},
	DisableFlagsInUseLine: true,
	Run: run(getOrganization, organizationTemplate),
}

var organizationCreateCmd = &cobra.Command{
	Use:                   "create [org] [billingEmail] [displayName]",
	Short:                 "Create a new organization",
	Args:                  cobra.ExactArgs(3),
	Aliases:               []string{"c"},
	DisableFlagsInUseLine: true,
	Run: run(createOrganization, organizationTemplate),
}

func init() {
	rootCmd.AddCommand(organizationCmd)

	organizationCmd.AddCommand(organizationGetCmd)
	organizationCmd.AddCommand(organizationCreateCmd)
}
