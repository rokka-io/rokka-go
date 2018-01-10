package cli

import (
	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

func createUser(c *rokka.Client, args []string) (interface{}, error) {
	return c.CreateUser(args[0], args[1])
}

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:                   "users",
	Short:                 "Manage users and organizations",
	Run:                   nil,
	Aliases:               []string{"u"},
	DisableFlagsInUseLine: true,
}

var usersCreateCmd = &cobra.Command{
	Use:                   "create [org] [email]",
	Short:                 "Create a new organization and a user with the given email address",
	Args:                  cobra.ExactArgs(1),
	Aliases:               []string{"c"},
	DisableFlagsInUseLine: true,
	Run: run(createUser, "Id:\t{{.ID}}\nEmail:\t{{.Email}}\nAPI Key:\t{{.APIKey}}\n\nAn email containing further information has been sent to you.\n"),
}

func init() {
	rootCmd.AddCommand(usersCmd)

	usersCmd.AddCommand(usersCreateCmd)
}
