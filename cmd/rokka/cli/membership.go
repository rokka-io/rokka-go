package cli

import (
	"fmt"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

// membershipRole contains a map of strings to roles.
var membershipRole = map[string]rokka.MembershipRole{
	"read":  rokka.RoleRead,
	"write": rokka.RoleWrite,
	"admin": rokka.RoleAdmin,
}

func createMembership(c *rokka.Client, args []string) (interface{}, error) {
	role, ok := membershipRole[args[2]]
	if !ok {
		return nil, fmt.Errorf(`Invalid membership role "%s"`, args[2])
	}
	return nil, c.CreateMembership(args[0], args[1], role)
}

// membershipCmd represents the membership command
var membershipCmd = &cobra.Command{
	Use:                   "membership",
	Short:                 "Manage memberships (assign new users to organizations, change roles, etc.)",
	Run:                   nil,
	Aliases:               []string{"m"},
	DisableFlagsInUseLine: true,
}

var membershipCreateCmd = &cobra.Command{
	Use:                   "create [org] [email] [role]",
	Short:                 "Create Membership of a user to an organization with the specified role (read, write, admin)",
	Args:                  cobra.ExactArgs(3),
	Aliases:               []string{"c"},
	DisableFlagsInUseLine: true,
	Run: run(createMembership, "Membership successfully created.\n"),
}

func init() {
	rootCmd.AddCommand(membershipCmd)

	membershipCmd.AddCommand(membershipCreateCmd)
}
