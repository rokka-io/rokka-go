package rokka

import (
	"fmt"
	"time"
)

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

type GetSourceImageResponse struct {
	Hash         string                 `json:"hash"`
	ShortHash    string                 `json:"short_hash"`
	BinaryHash   string                 `json:"binary_hash"`
	Created      time.Time              `json:"created"`
	Name         string                 `json:"name"`
	MimeType     string                 `json:"mimetype"`
	Format       string                 `json:"format"`
	Size         int                    `json:"size"`
	Width        int                    `json:"width"`
	Height       int                    `json:"height"`
	Organization string                 `json:"organization"`
	Link         string                 `json:"link"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
}

func (c *Client) ListSourceImages(org string, query map[string]string) (ListSourceImagesResponse, error) {
	result := ListSourceImagesResponse{}

	req, err := c.NewRequest("GET", "/sourceimages/"+org, nil, query)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)

	return result, err
}

func (c *Client) GetSourceImage(org, hash string) (GetSourceImageResponse, error) {
	result := GetSourceImageResponse{}

	req, err := c.NewRequest("GET", fmt.Sprintf("/sourceimages/%s/%s", org, hash), nil, nil)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)

	return result, err
}
