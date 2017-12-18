package rokka

import (
	"net/http"
	"time"
)

type statsResponseValue struct {
	Timestamp time.Time `json:"timestamp"`
	Value     int       `json:"value"`
	Unit      string    `json:"unit,omitempty"`
}

// StatsResponse contains time based statistics for an organization.
type StatsResponse struct {
	SpaceInBytes    []statsResponseValue `json:"space_in_bytes"`
	NumberOfFiles   []statsResponseValue `json:"number_of_files"`
	BytesDownloaded []statsResponseValue `json:"bytes_downloaded"`
}

// GetStats retrieves statistics for an organization.
//
// See: https://rokka.io/documentation/references/stats.html
func (c *Client) GetStats(org string, query map[string]string) (StatsResponse, error) {
	result := StatsResponse{}

	req, err := c.NewRequest(http.MethodGet, "/stats/"+org, nil, query)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)
	return result, err
}
