package cli

import (
	"errors"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

var (
	errInvalidAPIKey  = errors.New("invalid API key")
	errAPIKeyRequired = errors.New("missing flag --apiKey")
)

func login(c *rokka.Client, args []string) (interface{}, error) {
	valid, err := c.ValidAPIKey()
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errInvalidAPIKey
	}

	cfg := Config{
		APIKey:    c.GetConfig().APIKey,
		ImageHost: c.GetConfig().ImageHost,
	}
	err = SaveConfig(cfg)
	if err != nil {
		return nil, err
	}
	return "Login successful", nil
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Stores the API key as JSON in the configuration file specified using `--config`.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Lookup("apiKey").Changed {
			return errAPIKeyRequired
		}
		return nil
	},
	Run:     run(login, "{{.}}\n"),
	Aliases: []string{"l"},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
