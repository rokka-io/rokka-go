package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/rokka-io/rokka-go/cmd/rokka/cli/massupload"
	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

type uploadResult struct {
	Path      string
	RokkaHash string
	Error     error
}

type massuploadOptions struct {
	Organization  string
	BasePath      string
	Recursive     bool
	DryRun        bool
	Extensions    []string
	Concurrency   int
	Template      string
	TemplateError string
}

var massuploadCmd = &cobra.Command{
	Use:   "massupload [organization] [path]",
	Short: "Upload all images from a folder to rokka",
	Args:  cobra.ExactArgs(2),
	DisableFlagsInUseLine: true,
	Run: run(massUpload, "Uploaded images: {{.TotalUploads}}\n"),
}

var options massuploadOptions

func init() {
	rootCmd.AddCommand(massuploadCmd)

	flags := massuploadCmd.Flags()
	flags.BoolVarP(
		&options.Recursive,
		"recursive",
		"",
		false,
		"Recurse over the folder",
	)
	flags.StringSliceVarP(
		&options.Extensions,
		"extensions",
		"",
		[]string{"gif", "jpg", "png"},
		"Only upload the given file extensions --extensions=gif,jpg",
	)
	flags.IntVarP(
		&options.Concurrency,
		"concurrency",
		"",
		1,
		"Number of concurrent processes to use for uploading images",
	)
	flags.BoolVarP(
		&options.DryRun,
		"dry-run",
		"",
		false,
		"Simulate operation, do not upload files to Rokka.io",
	)
	flags.StringVarP(
		&options.Template,
		"template-success",
		"",
		"Uploaded {{.Path}} {{.RokkaHash}}\n",
		"Template to be applied to successful uploads (See: https://golang.org/pkg/text/template/)",
	)
	flags.StringVarP(
		&options.TemplateError,
		"template-error",
		"",
		"Failure for {{.Path}} {{.Error}}\n",
		"Template to be applied to error uploads (See: https://golang.org/pkg/text/template/)",
	)
}

// Validate the MassuploadOptions: returns an error, or nil if everything is OK
func (m *massuploadOptions) validate() error {
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

func massUpload(c *rokka.Client, args []string) (interface{}, error) {
	options.Organization = args[0]
	options.BasePath = args[1]

	validation := options.validate()

	if nil != validation {
		return nil, validation
	}

	images := make(chan string)
	results := make(chan uploadResult)

	startWorkers(options, rokkaClient, images, results)

	// Scan folders and files
	go scan(options.BasePath, images, options.Recursive, options.Extensions)

	totalUploads := 0
	totalFailures := 0

	okReporter, err := massupload.BuildReporter(options.Template, logger.StdOut)
	if err != nil {
		logErrorAndExit(err)
	}

	errorReporter, err := massupload.BuildReporter(options.TemplateError, logger.StdOut)
	if err != nil {
		logErrorAndExit(err)
	}

	// Collect results and display progress
	for result := range results {
		// Rendering via template the successful uploads, or failures
		if nil != result.Error {
			totalFailures += 1
			err := errorReporter.Report(result)
			if err != nil {
				logErrorAndExit(err)
			}
		} else {
			totalUploads += 1
			err := okReporter.Report(result)
			if err != nil {
				logErrorAndExit(err)
			}
		}
	}

	return struct {
		TotalUploads  int
		TotalFailures int
	}{TotalFailures: totalFailures, TotalUploads: totalUploads}, nil
}

func startWorkers(options massuploadOptions, client *rokka.Client, images chan string, results chan uploadResult) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(options.Concurrency)

	userMetadata := buildUserMetadata(options)

	// Start workers for image upload
	for i := 0; i < options.Concurrency; i++ {
		go func() {
			defer waitGroup.Done()
			uploadWorker(client, images, results, options.Organization, userMetadata, options.DryRun)
		}()
	}

	// Start a go-routine to to close the result channel as soon as all workers are done
	go func() {
		waitGroup.Wait()
		close(results)
	}()

}

func buildUserMetadata(options massuploadOptions) map[string]interface{} {
	// TODO: Build the user-metadata from the input options
	return nil
}

func scan(root string, images chan string, recursive bool, extensions []string) error {
	// Keep extensions sorted to use a binary search for matching
	sort.Strings(extensions)

	defer func() {
		close(images)
	}()

	err := filepath.Walk(massupload.Fixpath(root), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info == nil {
			return nil
		}

		// Skip subfolders if enabled, but still scan the root directory
		if info.IsDir() && root != path && !recursive {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			// Exclude files without extension
			ext := strings.TrimPrefix(filepath.Ext(path), ".")
			if ext == "" {
				return nil
			}

			// Exclude empty files
			if info.Size() == 0 {
				return nil
			}

			i := sort.SearchStrings(extensions, ext)
			if i >= len(extensions) || extensions[i] != ext {
				// If the current extension is not found among the allowed ones, just skip the file
				return nil
			}

			// Add the image to the list of images to be uploaded
			images <- path
		}

		return nil
	})

	return err
}

func uploadWorker(
	client *rokka.Client,
	imageFiles chan string,
	results chan uploadResult,
	organization string,
	userMetadata map[string]interface{},
	dryRun bool,
) {
	for path := range imageFiles {
		result := uploadResult{
			Path: path,
		}

		if dryRun {
			result.RokkaHash = "DRY-RUN-HASH"
		} else {
			response, err := executeRokkaUpload(client, path, organization, userMetadata)
			if err != nil {
				result.Error = err
			} else {
				result.RokkaHash = response.Items[0].Hash
			}
		}

		results <- result
	}
}

func executeRokkaUpload(
	client *rokka.Client,
	image string,
	organization string,
	userMetadata map[string]interface{},
) (rokka.CreateSourceImageResponse, error) {
	file, err := os.Open(image)
	if err != nil {
		return rokka.CreateSourceImageResponse{}, err
	}

	defer file.Close()

	return client.CreateSourceImageWithMetadata(organization, filepath.Base(image), file, userMetadata, nil)
}
