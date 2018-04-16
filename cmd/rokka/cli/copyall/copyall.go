package copyall

import (
	"errors"
	"sync"

	"github.com/rokka-io/rokka-go/rokka"
)

//Options general options for the whole copy-all command
type Options struct {
	SourceOrganization      string
	DestinationOrganization string
	DryRun                  bool
	Concurrency             int
}

type listReturn struct {
	Cursor     string
	ItemsCount int
}

//CopyResult the results of the copy operation
type CopyResult struct {
	RokkaHash string
	Error     error
}

//StartWorkers starts Copy Workers
func StartWorkers(options Options, client *rokka.Client, images chan string, results chan CopyResult, quit chan bool) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(options.Concurrency)

	// Start workers for image copy
	for i := 0; i < options.Concurrency; i++ {
		go func() {
			defer waitGroup.Done()
			copyWorker(client, images, results, quit, options.SourceOrganization, options.DestinationOrganization, options.DryRun)
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
	quit chan bool,
	sourceOrg string,
	destinationOrg string,
	dryRun bool,
) {
	for hash := range imageFiles {
		result := CopyResult{
			RokkaHash: hash,
		}

		if dryRun {
			result.RokkaHash = "DRY-RUN-HASH"
		} else {
			err := executeRokkaCopy(client, hash, sourceOrg, destinationOrg)
			if err != nil {
				result.Error = err
			}
		}

		select {
		case results <- result:
		case <-quit: // If we get a quit message, end the worker!
			return
		}
	}
}

//Scan starts the scan operation for getting all images
func Scan(options Options, client *rokka.Client, images chan string, quit chan bool) error {
	defer func() {
		close(images)
	}()
	doIt := true
	cursor := ""
	for doIt {
		listReturnValues, err := list(options, client, images, quit, cursor)
		if err != nil {
			return err
		}

		if listReturnValues.Cursor == "" || cursor == listReturnValues.Cursor || listReturnValues.ItemsCount == 0 {
			doIt = false
		} else {
			cursor = listReturnValues.Cursor
		}
	}
	return nil
}

func list(options Options, client *rokka.Client, images chan string, quit chan bool, cursor string) (listReturn, error) {
	listSourceImagesOptions := rokka.ListSourceImagesOptions{}
	listReturnValues := listReturn{"", 0}
	if cursor != "" {
		listSourceImagesOptions.Offset = cursor
	}
	res, err := client.ListSourceImages(options.SourceOrganization, listSourceImagesOptions)
	if err != nil {
		return listReturnValues, err
	}
	for _, element := range res.Items {
		// Add the image to the list of ones to be uploaded
		select {
		case images <- element.Hash:
		case <-quit:
			return listReturnValues, errors.New("image copy cancelled")
		}
	}
	listReturnValues.Cursor = res.Cursor
	listReturnValues.ItemsCount = len(res.Items)
	return listReturnValues, nil
}

func executeRokkaCopy(client *rokka.Client, hash string, sourceOrg string, destinationOrg string) error {
	return client.CopySourceImage(sourceOrg, hash, destinationOrg)
}
