package copyall

import (
	"sync"

	"github.com/rokka-io/rokka-go/rokka"
)

//Options general options for the whole copy-all command
type Options struct {
	SourceOrganization      string
	DestinationOrganization string
	DryRun                  bool
	Concurrency             int
	NoProgress              bool
}

// CopyResult contains the result of the operatio
type CopyResult struct {
	RokkaHash string
	Error     error
}

// StartWorkers starts Copy Workers
func StartWorkers(options Options, client *rokka.Client, images chan string, results chan CopyResult) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(options.Concurrency)

	// Start workers for image copy
	for i := 0; i < options.Concurrency; i++ {
		go func() {
			defer waitGroup.Done()
			copyWorker(client, images, results, options.SourceOrganization, options.DestinationOrganization, options.DryRun)
		}()
	}

	// Start a go-routine to to close the result channel as soon as all workers are done
	go func() {
		waitGroup.Wait()
		close(results)
	}()

}

func copyWorker(
	client *rokka.Client,
	imageFiles chan string,
	results chan CopyResult,
	sourceOrg string,
	destinationOrg string,
	dryRun bool,
) {
	for hash := range imageFiles {
		result := CopyResult{
			RokkaHash: hash,
		}

		if !dryRun {
			result.Error = executeRokkaCopy(client, hash, sourceOrg, destinationOrg)
		}

		results <- result

	}
}

// Scan starts the scan operation for getting all images
func Scan(options Options, client *rokka.Client, images chan string) error {
	defer close(images)
	cursor := ""
	for {
		newCursor, itemsCount, err := list(options, client, images, cursor)
		if err != nil {
			return err
		}

		if newCursor == "" || cursor == newCursor || itemsCount == 0 {
			return nil
		}
		cursor = newCursor
	}
}

func list(options Options, client *rokka.Client, images chan string, cursor string) (string, int, error) {
	listSourceImagesOptions := rokka.ListSourceImagesOptions{}
	if cursor != "" {
		listSourceImagesOptions.Offset = cursor
	}
	res, err := client.ListSourceImages(options.SourceOrganization, listSourceImagesOptions)
	if err != nil {
		return "", 0, err
	}
	for _, element := range res.Items {
		// Add the image to the list of ones to be copied
		images <- element.Hash
	}
	return res.Cursor, len(res.Items), err
}

func executeRokkaCopy(client *rokka.Client, hash string, sourceOrg string, destinationOrg string) error {
	return client.CopySourceImage(sourceOrg, hash, destinationOrg)
}
