package cli

import (
	"fmt"
	"os"

	"os/user"

	"path"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

var (
	// auto-set during build of CLI
	cliVersion = "unversioned"

	apiKey           string
	apiAddress       string
	raw              bool
	verbose          bool
	responseTemplate string
	configFile       string

	logger *cliLog
	cl     *rokka.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "rokka",
	Version: cliVersion,
	Short:   "A CLI for rokka, the image hosting service",
	// TODO: add something helpful as a long text.
	//Long: ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Verbose = verbose

		hc := newHTTPClient(logger)

		cl = rokka.NewClient(&rokka.Config{
			APIKey:     apiKey,
			APIAddress: apiAddress,
			HTTPClient: hc,
		})
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	logger = NewCLILog(verbose)

	p, err := getPath()
	if err != nil {
		fmt.Println("Unable to get config path: " + err.Error())
		os.Exit(1)
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", p, "Config file to store the API key")
	rootCmd.PersistentFlags().StringVar(&apiKey, "apiKey", "", "API key")
	rootCmd.PersistentFlags().StringVar(&apiAddress, "apiAddress", "", "API address")
	rootCmd.PersistentFlags().BoolVarP(&raw, "raw", "r", false, "Show raw HTTP response")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode")
	rootCmd.PersistentFlags().StringVar(&responseTemplate, "template", "", "Template to be applied to the response (See: https://golang.org/pkg/text/template/)")
}

func getPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(usr.HomeDir, ".rokka", "config"), nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	SetConfigPath(configFile)

	cfg, err := GetConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(apiKey) == 0 {
		apiKey = cfg.APIKey
	}
}
