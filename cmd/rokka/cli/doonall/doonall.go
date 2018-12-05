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

// OperationResult contains the result of the operation
type OperationResult struct {
	NotOK int
	OK    int
	Error error
}

// CallbackFunc is run on each image of an organization
type CallbackFunc func(client *rokka.Client, hashes []string, options Options) OperationResult

// StartWorkers starts Copy Workers
func StartWorkers(options Options, client *rokka.Client, images chan string, results chan OperationResult, callback CallbackFunc, limit int) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(options.Concurrency)
	// Start workers for image copy
	for i := 0; i < options.Concurrency; i++ {
		go func() {
			defer waitGroup.Done()
			hashCallbackWorker(client, images, results, options, callback, limit)
		}()
	}
	// Start a go-routine to to close the result channel as soon as all workers are done
	go func() {
		waitGroup.Wait()
		close(results)
	}()
}

func hashCallbackWorker(
	client *rokka.Client,
	imageFiles chan string,
	opresults chan OperationResult,
	options Options,
	callback CallbackFunc,
	limit int,
) {
	hashes := []string{}
	var opresult = OperationResult{}
	for hash := range imageFiles {
		hashes = append(hashes, hash)
		if len(hashes) >= limit {
			if !options.DryRun {
				opresult = callback(client, hashes, options)
			}

			hashes = []string{}
			opresults <- opresult
		}
	}
	if !options.DryRun {
		opresult = callback(client, hashes, options)
	}
	opresults <- opresult
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
func ExecuteRokkaCopy(client *rokka.Client, hashes []string, options Options) OperationResult {
	OK, notOK, err := client.CopySourceImages(options.SourceOrganization, hashes, options.DestinationOrganization)
	return OperationResult{OK: OK, NotOK: notOK, Error: err}

}

// ExecuteRokkaDelete delete one single image from the org
func ExecuteRokkaDelete(client *rokka.Client, hashes []string, options Options) OperationResult {
	var lasterr error

	OK := 0
	notOK := 0
	for _, hash := range hashes {
		lasterr = client.DeleteSourceImage(options.SourceOrganization, hash)
		if lasterr != nil {
			notOK++
		} else {
			OK++
		}
	}
	return OperationResult{OK: OK, NotOK: notOK, Error: lasterr}
}
