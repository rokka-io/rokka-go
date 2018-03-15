package batch

import (
	"github.com/rokka-io/rokka-go/rokka"
	"gopkg.in/cheggaaa/pb.v1"
)

// SourceImagesReader reads images from rokka on an organization.
type SourceImagesReader struct {
	Organization string
}

// Read uses the search API to paginate through all images.
func (sir *SourceImagesReader) Read(client *rokka.Client, images chan string, bar *pb.ProgressBar) error {
	cursor := ""
	for {
		opt := rokka.ListSourceImagesOptions{
			Offset: cursor,
		}
		res, err := client.ListSourceImages(sir.Organization, opt)
		if err != nil {
			return err
		}
		bar.Total = int64(res.Total)

		for _, element := range res.Items {
			images <- element.Hash
		}
		if res.Cursor == "" || cursor == res.Cursor || len(res.Items) == 0 {
			return nil
		}
		cursor = res.Cursor
	}
}

// Count fetches one image in order to get the Total amount of images available.
func (sir *SourceImagesReader) Count(client *rokka.Client) (int, error) {
	listSourceImagesOptions := rokka.ListSourceImagesOptions{
		Limit: 1,
	}
	res, err := client.ListSourceImages(sir.Organization, listSourceImagesOptions)
	if err != nil {
		return 0, err
	}
	return res.Total, nil
}

// CopyAllSourceImagesWriter uses the copy all API to transfer images from one organization to a destination organization.
type CopyAllSourceImagesWriter struct {
	SourceOrganization      string
	DestinationOrganization string
}

func (cas *CopyAllSourceImagesWriter) Write(client *rokka.Client, images []string) OperationResult {
	OK, notOK, err := client.CopySourceImages(cas.SourceOrganization, images, cas.DestinationOrganization)
	return OperationResult{OK: OK, NotOK: notOK, Error: err}
}

// DeleteAllSourceImagesWriter deletes source images of an organization.
type DeleteAllSourceImagesWriter struct {
	Organization string
}

func (das *DeleteAllSourceImagesWriter) Write(client *rokka.Client, images []string) OperationResult {
	var lastErr error

	OK := 0
	notOK := 0
	for _, hash := range images {
		lastErr = client.DeleteSourceImage(das.Organization, hash)
		if lastErr != nil {
			notOK++
		} else {
			OK++
		}
	}
	return OperationResult{OK: OK, NotOK: notOK, Error: lastErr}
}

// NoopWriter does not do anything. It is used for the dry run.
type NoopWriter struct{}

func (nw *NoopWriter) Write(client *rokka.Client, images []string) OperationResult {
	return OperationResult{OK: len(images)}
}
