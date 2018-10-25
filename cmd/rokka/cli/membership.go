package cli

import (
	"fmt"
	"strings"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

// membershipRole contains a map of strings to roles.
var membershipRole = map[string]rokka.MembershipRole{
	"read":   rokka.RoleRead,
	"write":  rokka.RoleWrite,
	"admin":  rokka.RoleAdmin,
	"upload": rokka.RoleUpload,
}

func createMembership(c *rokka.Client, args []string) (interface{}, error) {
	roles, err := getRoles(args[2])
	if err != nil {
		return nil, err
	}
	return nil, c.CreateMembership(args[0], args[1], roles)
}

func createNewMembershipWithCurrentUser(c *rokka.Client, args []string) (interface{}, error) {
	roles, err := getRoles(args[1])
	if err != nil {
		return nil, err
	}
	return c.CreateNewMembershipWithCurrentUser(args[0], roles)
}

func listMembership(c *rokka.Client, args []string) (interface{}, error) {
	if len(args) == 2 {
		return c.ListMembershipForUUID(args[0], args[1])

	}
	return c.ListMembership(args[0])
}

func deleteMembership(c *rokka.Client, args []string) (interface{}, error) {
	return nil, c.DeleteMembership(args[0], args[1])
}

func getRoles(rolesIn string) ([]rokka.MembershipRole, error) {
	splittedRoles := strings.Split(rolesIn, ",")

	var roles []rokka.MembershipRole
	for _, r := range splittedRoles {
		role, ok := membershipRole[r]
		if !ok {
			return nil, fmt.Errorf(`Invalid membership role "%s"`, r)
		}
		roles = append(roles, role)
	}
	return roles, nil
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
	Use:                   "create [org] [uuid] [role]",
	Short:                 "Create Membership of a user to an organization with the specified comma seperated roles (read, write, upload, admin)",
	Args:                  cobra.ExactArgs(3),
	Aliases:               []string{"c"},
	DisableFlagsInUseLine: true,
	Run:                   run(createMembership, "Membership successfully created."),
}
var membershipCreateWithUserCmd = &cobra.Command{
	Use:                   "createWithUser [org] [roles]",
	Short:                 "Create new Membership for the current user to an organization with the specified comma seperated  oles (read, write, upload, admin)",
	Args:                  cobra.ExactArgs(2),
	Aliases:               []string{"cu"},
	DisableFlagsInUseLine: true,
	Run:                   run(createNewMembershipWithCurrentUser, "Membership and user successfully created.\nEmail\tUUID\tRoles\tApi-Key\n{{.Email}}\t{{.UserID}}\t{{range $i, $e := .Roles}}{{if $i}},{{end}}{{$e}}{{end}}\t{{.APIKey}}\n"),
}

var membershipListCmd = &cobra.Command{
	Use:                   "list [org] ",
	Short:                 "List Memberships of an organization",
	Args:                  cobra.ExactArgs(1),
	Aliases:               []string{"l"},
	DisableFlagsInUseLine: true,
	Run:                   run(listMembership, "Email\tUUID\tRoles\n{{range .Items}}{{.Email}}\t{{.UserID}}\t{{range $i, $e := .Roles}}{{if $i}},{{end}}{{$e}}{{end}}\n{{end}}"),
}

var membershipListForUUIDCmd = &cobra.Command{
	Use:                   "get [org] [uuid]",
	Short:                 "Get Membership of a user of an organization",
	Args:                  cobra.ExactArgs(2),
	Aliases:               []string{"g"},
	DisableFlagsInUseLine: true,
	Run:                   run(listMembership, "Email\tUUID\tRoles\n{{.Email}}\t{{.UserID}}\t{{range $i, $e := .Roles}}{{if $i}},{{end}}{{$e}}{{end}}\n"),
}

var membershipDeleteCmd = &cobra.Command{
	Use:                   "delete [org] [uuid]",
	Short:                 "Delete Membership of a user of an organization",
	Args:                  cobra.ExactArgs(2),
	Aliases:               []string{"d"},
	DisableFlagsInUseLine: true,
	Run:                   run(deleteMembership, "Membership successfully deleted.\n"),
}

func init() {
	rootCmd.AddCommand(membershipCmd)

	membershipCmd.AddCommand(membershipCreateCmd)
	membershipCmd.AddCommand(membershipCreateWithUserCmd)
	membershipCmd.AddCommand(membershipListCmd)
	membershipCmd.AddCommand(membershipListForUUIDCmd)
	membershipCmd.AddCommand(membershipDeleteCmd)

}
