package batch

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rokka-io/rokka-go/rokka"
	"gopkg.in/cheggaaa/pb.v1"
)

// MassUploadOptions are specific CLI flags for the mass upload CLI cmd.
type MassUploadOptions struct {
	Recursive  bool
	Extensions []string
}

// MassUploader is both a Reader and Writer which reads from the fileSystem and creates source images in the writer.
type MassUploader struct {
	BasePath     string
	Recursive    bool
	Extensions   []string
	Organization string
	UserMetadata map[string]interface{}
}

// Read walks the directory specified in the CLI and adds the found images (filtered by extensions) to the image channel.
func (mu *MassUploader) Read(client *rokka.Client, images chan string, bar *pb.ProgressBar) error {
	// Keep extensions sorted to use a binary search for matching
	extensions := mu.Extensions
	sort.Strings(extensions)

	return filepath.Walk(Fixpath(mu.BasePath), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info == nil {
			return nil
		}

		// Skip subfolders if enabled, but still scan the root directory
		if info.IsDir() && mu.BasePath != path && !mu.Recursive {
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

			bar.Total++

			// Add the image to the list of images to be uploaded
			images <- path
		}

		return nil
	})
}

// Write creates source images for each image.
func (mu *MassUploader) Write(client *rokka.Client, images []string) OperationResult {
	var lastErr error

	OK := 0
	notOK := 0
	for _, path := range images {
		if err := mu.uploadFile(client, path); err != nil {
			lastErr = err
			notOK++
		} else {
			OK++
		}
	}
	return OperationResult{OK: OK, NotOK: notOK, Error: lastErr}
}

func (mu *MassUploader) uploadFile(client *rokka.Client, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = client.CreateSourceImageWithMetadata(mu.Organization, filepath.Base(path), file, mu.UserMetadata, nil)
	return err
}
