package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rokka-io/rokka-go/cmd/rokka/cli/batch"
	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	batchOptions      batch.Options
	massUploadOptions batch.MassUploadOptions
)

func copyAllSourceImage(c *rokka.Client, args []string) (interface{}, error) {
	sourceOrganization := args[0]
	destinationOrganization := args[1]

	sir := batch.SourceImagesReader{Organization: sourceOrganization}
	cas := batch.CopyAllSourceImagesWriter{SourceOrganization: sourceOrganization, DestinationOrganization: destinationOrganization}

	return executeBatchCmd(c, batchOptions, &cas, &sir, &sir, fmt.Sprintf("Copying of %%d source images from organization %s to %s\n", sourceOrganization, destinationOrganization), 100)
}

func deleteAllSourceImage(c *rokka.Client, args []string) (interface{}, error) {
	sourceOrganization := args[0]

	sir := batch.SourceImagesReader{Organization: sourceOrganization}
	das := batch.DeleteAllSourceImagesWriter{Organization: sourceOrganization}

	return executeBatchCmd(c, batchOptions, &das, &sir, &sir, fmt.Sprintf("Deleting of %%d source images on organization %s.\n", sourceOrganization), 1)
}

func massUpload(c *rokka.Client, args []string) (interface{}, error) {
	organization := args[0]
	basePath := args[1]

	mu := batch.MassUploader{
		BasePath:     basePath,
		Organization: organization,
		Recursive:    massUploadOptions.Recursive,
		Extensions:   massUploadOptions.Extensions,
	}

	return executeBatchCmd(c, batchOptions, &mu, &mu, nil, fmt.Sprintf("Uploading images from directory `%s` to organization `%s`.\n", basePath, organization), 100)
}

var sourceImagesCopyAllCmd = &cobra.Command{
	Use:                   "copy-all [sourceOrg] [destinationOrg]",
	Short:                 "Copy all source images from on org to another",
	Args:                  cobra.ExactArgs(2),
	Aliases:               []string{"cpa"},
	DisableFlagsInUseLine: true,
	Run:                   run(copyAllSourceImage, "Successfully copied {{.SuccessfullyUploaded}} source images. Errors with {{.ErrorUploaded}} source images.\n"),
}

var sourceImagesDeleteAllCmd = &cobra.Command{
	Use:                   "delete-all [org]",
	Short:                 "Deletes all source images from on org to another",
	Args:                  cobra.ExactArgs(1),
	Aliases:               []string{"del-all"},
	DisableFlagsInUseLine: true,
	Run:                   run(deleteAllSourceImage, "Successfully deleted {{.SuccessfullyUploaded}} source images. Errors with {{.ErrorUploaded}} source images.\n"),
}

var massUploadCmd = &cobra.Command{
	Use:   "massupload [organization] [path]",
	Short: "Upload all images from a folder to rokka",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("organization and path are required")
		}

		// Check if the path exists
		fileInfo, err := os.Stat(args[1])
		if os.IsNotExist(err) {
			return fmt.Errorf("path '%s' does not exist", args[1])
		}

		// Check if the path is a folder
		if !fileInfo.IsDir() {
			return fmt.Errorf("path '%s' must be a directory", args[1])
		}

		return nil
	},
	DisableFlagsInUseLine: true,
	Run:                   run(massUpload, "Successfully uploaded {{.SuccessfullyUploaded}} images. Errors with {{.ErrorUploaded}} images.\n"),
}

func executeBatchCmd(c *rokka.Client, options batch.Options, w batch.Writer, r batch.Reader, p batch.ProgressCounter, tmpl string, limit int) (interface{}, error) {
	images := make(chan string)
	results := make(chan batch.OperationResult)

	if options.DryRun {
		w = &batch.NoopWriter{}
	}

	go batch.WriteImages(c, images, results, w, options.Concurrency, limit)

	counterError, counterSuccess := 0, 0

	var total int
	var err error
	if p != nil {
		total, err = p.Count(c)
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(tmpl, "%d") {
		logger.Errorf(tmpl, total)
	} else {
		logger.Error(tmpl)
	}

	if !options.Force {
		logger.Errorf("Are you sure? (yes/no): ")
		if !askForConfirmation() {
			return nil, errors.New("operation cancelled")
		}
	}
	bar := pb.New(total)
	bar.ShowSpeed = true
	bar.Output = logger.StdErr
	if batchOptions.NoProgress {
		bar.NotPrint = true
	}

	go func() {
		if err := batch.ReadImages(c, images, r, bar); err != nil {
			logger.Errorf("Error reading images: %s", err)
		}
	}()

	bar.Start()
	for result := range results {
		counterSuccess += result.OK
		counterError += result.NotOK
		bar.Set(counterError + counterSuccess)

		if result.Error != nil {
			logger.Errorf("error writing: %s\n", result.Error)
		}
	}
	bar.Finish()

	return struct {
		SuccessfullyUploaded int
		ErrorUploaded        int
	}{counterSuccess, counterError}, nil
}

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		logger.Errorf("%s", err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if ListContains(okayResponses, response) {
		return true
	} else if ListContains(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

func addBatchFlags(f *flag.FlagSet) {
	f.IntVarP(
		&batchOptions.Concurrency,
		"concurrency",
		"",
		2,
		"Number of concurrent processes to use for uploading images",
	)
	f.BoolVarP(
		&batchOptions.DryRun,
		"dry-run",
		"",
		false,
		"Simulate operation, do not copy files on Rokka.io",
	)
	f.BoolVarP(
		&batchOptions.NoProgress,
		"no-progress",
		"",
		false,
		"No progress bar",
	)
	f.BoolVarP(
		&batchOptions.Force,
		"force",
		"f",
		false,
		"Don't request confirmation for executing the command",
	)
}

func init() {
	sourceImagesCmd.AddCommand(sourceImagesCopyAllCmd)
	sourceImagesCmd.AddCommand(sourceImagesDeleteAllCmd)
	sourceImagesCmd.AddCommand(massUploadCmd)

	addBatchFlags(sourceImagesCopyAllCmd.Flags())
	addBatchFlags(sourceImagesDeleteAllCmd.Flags())
	addBatchFlags(massUploadCmd.Flags())

	massUploadCmd.Flags().BoolVarP(
		&massUploadOptions.Recursive,
		"recursive",
		"",
		false,
		"Recurse over the folder",
	)
	massUploadCmd.Flags().StringSliceVarP(
		&massUploadOptions.Extensions,
		"extensions",
		"e",
		[]string{"gif", "jpg", "png"},
		"Only upload the given file extensions --extensions=gif,jpg",
	)
}
