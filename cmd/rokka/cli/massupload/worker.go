package massupload

import (
	"github.com/rokka-io/rokka-go/rokka"
	"os"
	"path/filepath"
)

func UploadWorker(
	client *rokka.Client,
	imageFiles chan string,
	results chan UploadResult,
	quit chan bool,
	organization string,
	userMetadata map[string]interface{},
	dryRun bool,
) {
	for path := range imageFiles {
		result := UploadResult{
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

		select {
		case results <- result:
		case <-quit: // If we get a quit message, end the worker!
			return
		}
	}
}

func executeRokkaUpload(client *rokka.Client, image string, organization string, userMetadata map[string]interface{}) (rokka.CreateSourceImageResponse, error) {
	file, err := os.Open(image)
	if err != nil {
		return rokka.CreateSourceImageResponse{}, err
	}

	// Postpone the file closing
	defer file.Close()

	return client.CreateSourceImageWithMetadata(organization, filepath.Base(image), file, userMetadata, nil)
}
