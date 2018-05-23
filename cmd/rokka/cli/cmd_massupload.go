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
	"sort"
	"strings"
)

type uploadResult struct {
	Path      string
	RokkaHash string
	Error     error
}

type imageDetails struct {
	Organization string
	Path         string
	UserMetadata map[string]interface{}
}

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
	results := make(chan uploadResult)

	startWorkers(options, rokkaClient, images, results)

	// Scan folders and files
	go scan(options.BasePath, images, options.Recursive, options.Extensions)

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

func startWorkers(options MassuploadOptions, client *rokka.Client, images chan string, results chan uploadResult) {
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

func buildUserMetadata(options MassuploadOptions) map[string]interface{} {
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

			fmt.Printf("Scan publish image %s\n", path)
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
