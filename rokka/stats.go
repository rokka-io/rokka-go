package rokka

import "time"

type StatsResponseValue struct {
	Timestamp time.Time `json:"timestamp"`
	Value     int       `json:"value"`
	Unit      string    `json:"unit,omitempty"`
}

type StatsResponse struct {
	SpaceInBytes    []StatsResponseValue `json:"space_in_bytes"`
	NumberOfFiles   []StatsResponseValue `json:"number_of_files"`
	BytesDownloaded []StatsResponseValue `json:"bytes_downloaded"`
}

func (c *Client) GetStats(name string, query map[string]string) (StatsResponse, error) {
	result := StatsResponse{}

	req, err := c.NewRequest("GET", "/stats/"+name, nil, query)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)
	return result, err
}
