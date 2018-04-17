package doonall

import (
	"sync"

	"github.com/rokka-io/rokka-go/rokka"
)

// Options general options for the whole copy-all command
type Options struct {
	SourceOrganization      string
	DestinationOrganization string
	DryRun                  bool
	Concurrency             int
	NoProgress              bool
	Force                   bool
}

// CopyResult contains the result of the operatio
type CopyResult struct {
	RokkaHash string
	Error     error
}

// CallbackFunc is run on each image of an organization
type CallbackFunc func(client *rokka.Client, hash string, options Options) (err error)

// StartWorkers starts Copy Workers
func StartWorkers(options Options, client *rokka.Client, images chan string, results chan CopyResult, callback CallbackFunc) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(options.Concurrency)
	// Start workers for image copy
	for i := 0; i < options.Concurrency; i++ {
		go func() {
			defer waitGroup.Done()
			copyWorker(client, images, results, options, callback)
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
	options Options,
	callback CallbackFunc) {
	for hash := range imageFiles {
		result := CopyResult{
			RokkaHash: hash,
		}
		if !options.DryRun {
			result.Error = callback(client, hash, options)
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

// ExecuteRokkaCopy copies a single image from one org to another
func ExecuteRokkaCopy(client *rokka.Client, hash string, options Options) error {
	return client.CopySourceImage(options.SourceOrganization, hash, options.DestinationOrganization)
}

// ExecuteRokkaDelete delete one single image from the org
func ExecuteRokkaDelete(client *rokka.Client, hash string, options Options) error {
	return client.DeleteSourceImage(options.SourceOrganization, hash)
}
