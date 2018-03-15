package batch

import (
	"sync"

	"github.com/rokka-io/rokka-go/rokka"
	"gopkg.in/cheggaaa/pb.v1"
)

// Options are used for the CLI flags for all batch commands
type Options struct {
	DryRun      bool
	Concurrency int
	NoProgress  bool
	Force       bool
}

// OperationResult contains the result of the operation
type OperationResult struct {
	NotOK int
	OK    int
	Error error
}

// Reader allows to read from an arbitrary location and inserts the image identifications to the channel for concurrent processing.
// The values within the images channel are string, and therefore can be pretty much anything. The writer used after reading
// needs to know and understand what an image string actually represents. E.g. in the case of copy-all it's an existing
// image hash. In the case of massUpload it's a filesystem path to an image.
type Reader interface {
	Read(client *rokka.Client, images chan string, bar *pb.ProgressBar) error
}

// Writer operates on the previously found image list and executes a write operation. This can be pretty much anything.
// For example it could create source images on rokka. Or delete a source image on rokka.
type Writer interface {
	Write(client *rokka.Client, images []string) OperationResult
}

// ProgressCounter returns the total images to be processed. It is used for the progress bar and confirmation messages of the batch
// CLIs. It's an optional interface to implement because not all operations allow to know beforehand how many images there are.
type ProgressCounter interface {
	Count(client *rokka.Client) (int, error)
}

// WriteImages creates a group of goroutines bound by the concurrency option. It executes the Writer.Write command for each flushInterval
// amount of images.
func WriteImages(client *rokka.Client, images chan string, results chan OperationResult, w Writer, concurrency int, flushInterval int) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(concurrency)
	defer close(results)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer waitGroup.Done()

			fileNames := make([]string, 0)

			for fileName := range images {
				fileNames = append(fileNames, fileName)

				if len(fileNames) >= flushInterval {
					results <- w.Write(client, fileNames)
					fileNames = make([]string, 0)
				}
			}

			// flush remaining items
			if len(fileNames) > 0 {
				results <- w.Write(client, fileNames)
			}
		}()
	}
	waitGroup.Wait()
}

// ReadImages is a simple wrapper around the Reader.Read call. In the future there may be things we can move here.
func ReadImages(client *rokka.Client, images chan string, r Reader, bar *pb.ProgressBar) error {
	defer close(images)
	return r.Read(client, images, bar)
}
