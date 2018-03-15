package cli

import (
	"errors"
	"fmt"
	"github.com/rokka-io/rokka-go/cmd/rokka/cli/massupload"
	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sync"
	// TODO: "syscall"
	// TOOD: "os/signal"
)

type MassuploadOptions struct {
	Organization string
	BasePath     string
	Recursive    bool
	DryRun       bool
	Extensions   []string
	Concurrency  int
}

var massuploadCmd = &cobra.Command{
	Use:   "massupload [organization] [path]",
	Short: "Upload all images from a folder to rokka",
	Args:  cobra.ExactArgs(2),
	DisableFlagsInUseLine: true,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		massuploadOptions.Organization = args[0]
		massuploadOptions.BasePath = args[1]

		return massuploadOptions.validate()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMassUpload(rokkaClient, massuploadOptions)
	},
}

var massuploadOptions MassuploadOptions

func init() {
	rootCmd.AddCommand(massuploadCmd)

	flags := massuploadCmd.Flags()
	flags.BoolVarP(
		&massuploadOptions.Recursive,
		"recursive",
		"",
		false,
		"Recurse over the folder",
	)
	flags.StringSliceVarP(
		&massuploadOptions.Extensions,
		"extensions",
		"",
		[]string{"gif", "jpg", "png"},
		"Only upload the given file extensions --extensions=gif,jpg",
	)
	flags.IntVarP(
		&massuploadOptions.Concurrency,
		"concurrency",
		"",
		1,
		"Number of concurrent processes to use for uploading images",
	)
	flags.BoolVarP(
		&massuploadOptions.DryRun,
		"dry-run",
		"",
		false,
		"Simulate operation, do not upload files to Rokka.io",
	)
}

// Validate the MassuploadOptions: returns an error, or nil if everything is OK
func (m *MassuploadOptions) validate() error {
	if m.Organization == "" {
		return errors.New("the 'organization' can not be empty")
	}

	if m.BasePath == "" {
		return errors.New("the 'path' can not be empty")
	}

	// Check if the path exists
	fileInfo, err := os.Stat(m.BasePath)
	if os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("path '%s' does not exist", m.BasePath))
	}

	// Check if the path is a folder
	if !fileInfo.IsDir() {
		return errors.New(fmt.Sprintf("path '%s' must be a directory", m.BasePath))
	}

	// Check if concurrency is valid
	if m.Concurrency < 1 {
		return errors.New(fmt.Sprint("concurrency option can not be less than 1!"))
	}

	return nil
}

func runMassUpload(client *rokka.Client, options MassuploadOptions) error {
	logger.Printf("Starting!\n%#v\n", options)

	images := make(chan string)
	results := make(chan massupload.UploadResult)
	quit := make(chan bool)

	// TODO: Handle os.Signals and send items into the quit channel!
	// signals := make(chan os.Signal, 1)
	// signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	//   sig := <-signals
	// 	 logger.Printf("Received signal %s, terminating workers", sig.String())
	// 	 close(quit)
	// 	 return
	// }()

	startWorkers(options, rokkaClient, images, results, quit)

	// Scan folders and files
	go massupload.Scan(options.BasePath, images, quit, options.Recursive, options.Extensions)

	// Collect results and display progress
	for result := range results {
		if result.Error != nil {
			logger.Errorf("Upload failed for %s! %s\n", filepath.Base(result.Path), result.Error)
		} else {
			logger.Printf("Uploaded %s, hash=%s\n", filepath.Base(result.Path), result.RokkaHash)
		}
	}

	return nil
}

func startWorkers(options MassuploadOptions, client *rokka.Client, images chan string, results chan massupload.UploadResult, quit chan bool) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(options.Concurrency)

	userMetadata := buildUserMetadata(options)

	// Start workers for image upload
	for i := 0; i < options.Concurrency; i++ {
		go func() {
			defer waitGroup.Done()
			massupload.UploadWorker(client, images, results, quit, options.Organization, userMetadata, options.DryRun)
		}()
	}

	// Start a go-routine to to close the result channel as soon as all workers are done
	go func() {
		waitGroup.Wait()
		close(results)
	}()

}

func buildUserMetadata(options MassuploadOptions) map[string]interface{} {
	// TODO: Build the user-metadata from the input options
	return nil
}
