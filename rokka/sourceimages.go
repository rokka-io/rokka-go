package rokka

import "time"

type ListSourceImagesResponse struct {
	Total int `json:"total"`
	Items []struct {
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
	} `json:"items"`
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
