package massupload

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Scan(root string, images chan string, quit chan bool, recursive bool, extensions []string) error {
	// Keep extensions sorted to use a binary search for matching
	sort.Strings(extensions)

	defer func() {
		close(images)
	}()

	err := filepath.Walk(fixpath(root), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info == nil {
			return nil
		}

		if info.IsDir() && !recursive {
			return filepath.SkipDir
		}

		if !info.IsDir() {

			// Exclude files without extension
			ext := strings.TrimPrefix(filepath.Ext(path), ".")
			if ext == "" {
				return nil
			}

			i := sort.SearchStrings(extensions, ext)
			if i >= len(extensions) || extensions[i] != ext {
				// If the current extension is not found among the allowed ones, just skip the file
				return nil
			}

			fmt.Printf("Scan publish image %s\n", path)
			// Add the image to the list of ones to be uploaded
			select {
			case images <- path:
			case <-quit:
				return errors.New("image scan cancelled")
			}
		}

		return nil
	})

	return err
}
