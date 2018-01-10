package rokka

import (
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

// GetStatsOptions defines the accepted querystring params for GetStats.
// Giving an empty struct will result in no querystring params being sent to rokka.
type GetStatsOptions struct {
	From string `url:"from,omitempty"`
	To   string `url:"to,omitempty"`
}

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
func (c *Client) GetStats(org string, options GetStatsOptions) (StatsResponse, error) {
	result := StatsResponse{}

	qs, err := query.Values(options)
	if err != nil {
		return result, err
	}

	req, err := c.NewRequest(http.MethodGet, "/stats/"+org, nil, qs)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)
	return result, err
}
