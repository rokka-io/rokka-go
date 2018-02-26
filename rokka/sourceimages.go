package rokka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/google/go-querystring/query"
)

// ListSourceImagesOptions defines the accepted query string params for GetStats.
// Giving an empty struct will result in no query string params being sent to rokka.
//
// See: https://rokka.io/documentation/references/searching-images.html#range-filtering
type ListSourceImagesOptions struct {
	Limit      int    `url:"limit,omitempty"`
	Offset     int    `url:"offset,omitempty"`
	Hash       string `url:"hash,omitempty"`
	BinaryHash string `url:"binaryhash,omitempty"`
	// Size can be an int or a range. See: https://github.com/rokka-io/rokka-go/issues/32
	Size   string `url:"size,omitempty"`
	Format string `url:"format,omitempty"`
	// Width can be an int or a range. See: https://github.com/rokka-io/rokka-go/issues/32
	Width string `url:"width,omitempty"`
	// Height can be an int or a range. See: https://github.com/rokka-io/rokka-go/issues/32
	Height string `url:"height,omitempty"`
	// Created needs to be always passed as a range.
	Created string `url:"created,omitempty"`
	// Sort by a specific field
	Sort string `url:"sort,omitempty"`
}

// ListSourceImagesResponse contains a list of source images alongside a total and pagination links.
type ListSourceImagesResponse struct {
	Total  int                      `json:"total"`
	Items  []GetSourceImageResponse `json:"items"`
	Cursor string                   `json:"cursor"`
	Links  struct {
		Prev *struct {
			Href string `json:"href,omitempty"`
		} `json:"prev,omitempty"`
		Next *struct {
			Href string `json:"href,omitempty"`
		} `json:"next,omitempty"`
	} `json:"links,omitempty"`
}

// GetSourceImageResponse is an object identifying an image.
type GetSourceImageResponse struct {
	Hash            string                 `json:"hash"`
	ShortHash       string                 `json:"short_hash"`
	BinaryHash      string                 `json:"binary_hash"`
	Created         time.Time              `json:"created"`
	Name            string                 `json:"name"`
	MimeType        string                 `json:"mimetype"`
	Format          string                 `json:"format"`
	Size            int                    `json:"size"`
	Width           int                    `json:"width"`
	Height          int                    `json:"height"`
	Organization    string                 `json:"organization"`
	Link            string                 `json:"link"`
	UserMetadata    map[string]interface{} `json:"user_metadata,omitempty"`
	DynamicMetadata map[string]interface{} `json:"dynamic_metadata,omitempty"`
}

// DownloadSourceImageResponse contains the original data as a byte slice and the filename.
type DownloadSourceImageResponse struct {
	Data     io.ReadCloser
	FileName string
}

// CreateSourceImageResponse is returned when creating an image.
type CreateSourceImageResponse struct {
	Total int                      `json:"total"`
	Items []GetSourceImageResponse `json:"items"`
}

// DynamicMetadataOptions defines the accepted options for adding dynamic metadata to an image.
type DynamicMetadataOptions struct {
	DeletePrevious bool `url:"deletePrevious,omitempty"`
}

// DynamicMetadataResponse contains the location of the updated image.
type DynamicMetadataResponse struct {
	Location string
}

var contentDispositionFilename *regexp.Regexp

func init() {
	contentDispositionFilename = regexp.MustCompile(`filename="([^"]+)"$`)
}

// ListSourceImages gets a paginated list of source images.
//
// See: https://rokka.io/documentation/references/searching-images.html
func (c *Client) ListSourceImages(org string, options ListSourceImagesOptions) (ListSourceImagesResponse, error) {
	result := ListSourceImagesResponse{}

	qs, err := query.Values(options)
	if err != nil {
		return result, err
	}

	req, err := c.NewRequest(http.MethodGet, "/sourceimages/"+org, nil, qs)
	if err != nil {
		return result, err
	}

	err = c.CallJSONResponse(req, &result)

	return result, err
}

// GetSourceImage returns the metadata of a single source image identified by its hash.
//
// See: https://rokka.io/documentation/references/source-images.html#retrieve-data-about-a-source-image
func (c *Client) GetSourceImage(org, hash string) (GetSourceImageResponse, error) {
	result := GetSourceImageResponse{}

	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("/sourceimages/%s/%s", org, hash), nil, nil)
	if err != nil {
		return result, err
	}

	err = c.CallJSONResponse(req, &result)

	return result, err
}

func downloadResponseHandler(resp *http.Response, v interface{}) error {
	if resp.StatusCode == http.StatusOK {
		v := v.(*DownloadSourceImageResponse)
		v.Data = resp.Body

		matches := contentDispositionFilename.FindStringSubmatch(resp.Header.Get("Content-Disposition"))
		v.FileName = matches[1]

		return nil
	}

	return handleStatusCodeError(resp)
}

// DownloadSourceImage allows to download the source image once uploaded.
//
// See: https://rokka.io/documentation/references/source-images.html
func (c *Client) DownloadSourceImage(org, hash string) (DownloadSourceImageResponse, error) {
	result := DownloadSourceImageResponse{}

	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("/sourceimages/%s/%s/download", org, hash), nil, nil)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result, downloadResponseHandler)

	return result, err
}

// DeleteSourceImage removes a source image by hash.
//
// See: https://rokka.io/documentation/references/source-images.html
func (c *Client) DeleteSourceImage(org, hash string) error {
	req, err := c.NewRequest(http.MethodDelete, fmt.Sprintf("/sourceimages/%s/%s", org, hash), nil, nil)
	if err != nil {
		return err
	}
	return c.Call(req, nil, nil)
}

// DeleteSourceImageByBinaryHash removes a source image by binaryhash.
//
// See: https://rokka.io/documentation/references/source-images.html
func (c *Client) DeleteSourceImageByBinaryHash(org, binaryHash string) error {
	qs := url.Values{}
	qs.Set("binaryHash", binaryHash)

	req, err := c.NewRequest(http.MethodDelete, fmt.Sprintf("/sourceimages/%s", org), nil, qs)
	if err != nil {
		return err
	}
	return c.Call(req, nil, nil)
}

// CreateSourceImage uploads an image without user or dynamic metadata set.
//
// See: https://rokka.io/documentation/references/source-images.html#create-a-source-image
func (c *Client) CreateSourceImage(org, name string, data io.Reader) (CreateSourceImageResponse, error) {
	return c.CreateSourceImageWithMetadata(org, name, data, nil, nil)
}

// CreateSourceImageWithMetadata uploads an image.
//
// See: https://rokka.io/documentation/references/source-images.html#create-a-source-image
func (c *Client) CreateSourceImageWithMetadata(org, name string, data io.Reader, userMetadata, dynamicMetadata map[string]interface{}) (CreateSourceImageResponse, error) {
	result := CreateSourceImageResponse{}

	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)
	fw, err := w.CreateFormFile("filename", name)
	if err != nil {
		return result, err
	}
	if _, err = io.Copy(fw, data); err != nil {
		return result, err
	}
	if userMetadata != nil {
		ffw, err := w.CreateFormField("meta_user[0]")
		if err != nil {
			return result, err
		}
		err = json.NewEncoder(ffw).Encode(userMetadata)
		if err != nil {
			return result, err
		}
	}
	for k, v := range dynamicMetadata {
		ffw, err := w.CreateFormField(fmt.Sprintf("meta_dynamic[0][%s]", k))
		if err != nil {
			return result, err
		}
		err = json.NewEncoder(ffw).Encode(v)
		if err != nil {
			return result, err
		}
	}
	w.Close()

	req, err := c.NewRequest(http.MethodPost, fmt.Sprintf("/sourceimages/%s", org), b, nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", w.FormDataContentType())
	err = c.CallJSONResponse(req, &result)

	return result, err
}

// dynamicMetadataResponseHandler is a responseHandler reading the Location header from the successful response.
func dynamicMetadataResponseHandler(resp *http.Response, v interface{}) error {
	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusCreated {
		v := v.(*DynamicMetadataResponse)
		v.Location = resp.Header.Get("Location")
		return nil
	}

	return handleStatusCodeError(resp)
}

// AddDynamicMetadata updates a source image by adding arbitrary metadata.
// Rokka generates a new image hash when calling this function. The return value of this call contains the location of the new image.
// If deletePrevious is true, the previous image will be deleted.
//
// See: https://rokka.io/documentation/references/dynamic-metadata.html
func (c *Client) AddDynamicMetadata(org, hash, name string, data io.Reader, options DynamicMetadataOptions) (DynamicMetadataResponse, error) {
	result := DynamicMetadataResponse{}

	qs, err := query.Values(options)
	if err != nil {
		return result, err
	}

	req, err := c.NewRequest(http.MethodPut, fmt.Sprintf("/sourceimages/%s/%s/meta/dynamic/%s", org, hash, name), data, qs)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result, dynamicMetadataResponseHandler)
	return result, err
}

// DeleteDynamicMetadata updates a source image by deleting existing metadata.
// Rokka generates a new image hash when calling this function. The return value of this call contains the location of the new image.
// If deletePrevious is true, the previous image will be deleted.
//
// See: https://rokka.io/documentation/references/dynamic-metadata.html
func (c *Client) DeleteDynamicMetadata(org, hash, name string, options DynamicMetadataOptions) (DynamicMetadataResponse, error) {
	result := DynamicMetadataResponse{}

	qs, err := query.Values(options)
	if err != nil {
		return result, err
	}

	req, err := c.NewRequest(http.MethodDelete, fmt.Sprintf("/sourceimages/%s/%s/meta/dynamic/%s", org, hash, name), nil, qs)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result, dynamicMetadataResponseHandler)
	return result, err
}

func (c *Client) userMetadata(method, org, hash string, data io.Reader) error {
	req, err := c.NewRequest(method, fmt.Sprintf("/sourceimages/%s/%s/meta/user", org, hash), data, nil)
	if err != nil {
		return err
	}

	return c.Call(req, nil, nil)
}

func (c *Client) userMetadataByName(method, org, hash, name string, data io.Reader) error {
	req, err := c.NewRequest(method, fmt.Sprintf("/sourceimages/%s/%s/meta/user/%s", org, hash, name), data, nil)
	if err != nil {
		return err
	}

	return c.Call(req, nil, nil)
}

// SetUserMetadata updates a source image by adding arbitrary metadata. If there were previous user metadata set on this source image,
// they'll get overwritten.
//
// See: https://rokka.io/documentation/references/image-metadata.html
func (c *Client) SetUserMetadata(org, hash string, data io.Reader) error {
	return c.userMetadata(http.MethodPut, org, hash, data)
}

// UpdateUserMetadata updates a source image by adding arbitrary metadata. Previous values get merged with the new values.
//
// See: https://rokka.io/documentation/references/image-metadata.html
func (c *Client) UpdateUserMetadata(org, hash string, data io.Reader) error {
	return c.userMetadata(http.MethodPatch, org, hash, data)
}

// DeleteUserMetadata updates a source image by deleting existing metadata.
//
// See: https://rokka.io/documentation/references/image-metadata.html
func (c *Client) DeleteUserMetadata(org, hash string) error {
	return c.userMetadata(http.MethodDelete, org, hash, nil)
}

// UpdateUserMetadataByName updates a source image by setting user metadata identified by name.
//
// See: https://rokka.io/documentation/references/image-metadata.html
func (c *Client) UpdateUserMetadataByName(org, hash, name string, data io.Reader) error {
	return c.userMetadataByName(http.MethodPut, org, hash, name, data)
}

// DeleteUserMetadataByName updates a source image by removing user metadata identified by name.
//
// See: https://rokka.io/documentation/references/image-metadata.html
func (c *Client) DeleteUserMetadataByName(org, hash, name string) error {
	return c.userMetadataByName(http.MethodDelete, org, hash, name, nil)
}
