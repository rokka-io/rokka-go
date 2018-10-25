package cli

import (
	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

func createUser(c *rokka.Client, args []string) (interface{}, error) {
	if len(args) == 2 {
		return c.CreateUser(args[0], args[1])

	}
	return c.CreateUserWithoutOrg(args[0])
}

func getUser(c *rokka.Client, args []string) (interface{}, error) {
	return c.GetUserID()
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
	Use:                   "create [email] [org] ",
	Short:                 "Create a new  user with the given email address and optional a new organization",
	Args:                  cobra.RangeArgs(1, 2),
	Aliases:               []string{"c"},
	DisableFlagsInUseLine: true,
	Run:                   run(createUser, "Id:\t{{.ID}}\nEmail:\t{{.Email}}\nAPI Key:\t{{.APIKey}}\n\nAn email containing further information has been sent to you.\n"),
}

var usersGetCmd = &cobra.Command{
	Use:                   "get ",
	Short:                 "Gets the UUID of the current user",
	Aliases:               []string{"g"},
	DisableFlagsInUseLine: true,
	Run:                   run(getUser, "Id:\t{{.ID}}\n"),
}

func init() {
	rootCmd.AddCommand(usersCmd)

	usersCmd.AddCommand(usersCreateCmd)
	usersCmd.AddCommand(usersGetCmd)

}
