package rokka

import "net/http"

type OrganizationResponse struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	Name         string `json:"name"`
	BillingEmail string `json:"billing_email"`
	Limit        struct {
		SpaceInBytes   *int `json:"space_in_bytes"`
		TrafficInBytes *int `json:"traffic_in_bytes"`
	} `json:"limit"`
}

func (c *Client) GetOrganization(name string) (OrganizationResponse, error) {
	result := OrganizationResponse{}

	req, err := c.NewRequest(http.MethodGet, "/organizations/"+name, nil, nil)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)
	return result, err
}
