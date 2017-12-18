package rokka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

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

type CreateSourceImageResponse struct {
	Total string                   `json:"total"`
	Items []GetSourceImageResponse `json:"items"`
}

// ListSourceImages gets a paginated list of source images.
//
// See: https://rokka.io/documentation/references/searching-images.html
func (c *Client) ListSourceImages(org string, query map[string]string) (ListSourceImagesResponse, error) {
	result := ListSourceImagesResponse{}

	req, err := c.NewRequest(http.MethodGet, "/sourceimages/"+org, nil, query)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)

	return result, err
}

// GetSourceImage returns the metadata of a single source image identified by it's hash.
//
// See: https://rokka.io/documentation/references/source-images.html#retrieve-data-about-a-source-image
func (c *Client) GetSourceImage(org, hash string) (GetSourceImageResponse, error) {
	result := GetSourceImageResponse{}

	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("/sourceimages/%s/%s", org, hash), nil, nil)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)

	return result, err
}

func (c *Client) CreateSourceImage(org, name string, data io.Reader) (CreateSourceImageResponse, error) {
	return c.CreateSourceImageWithMetadata(org, name, data, nil, nil)
}

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
	if dynamicMetadata != nil {
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
	}
	w.Close()

	req, err := c.NewRequest(http.MethodPost, fmt.Sprintf("/sourceimages/%s", org), b, nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", w.FormDataContentType())
	err = c.Call(req, &result)

	return result, err
}
