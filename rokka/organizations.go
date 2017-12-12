package rokka

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type OrganizationResponse struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	Name         string `json:"name"`
	BillingEmail string `json:"billing_email"`
	Limit        struct {
		SpaceInBytes   int `json:"space_in_bytes,omitempty"`
		TrafficInBytes int `json:"traffic_in_bytes,omitempty"`
	} `json:"limit,omitempty"`
}

type CreateOrganizationRequest struct {
	DisplayName  string `json:"display_name,omitempty"`
	BillingEmail string `json:"billing_email"`
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

func (c *Client) CreateOrganization(name, billingEmail, displayName string) (OrganizationResponse, error) {
	result := OrganizationResponse{}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(CreateOrganizationRequest{
		DisplayName:  displayName,
		BillingEmail: billingEmail,
	})
	if err != nil {
		return result, err
	}

	req, err := c.NewRequest(http.MethodPut, "/organizations/"+name, b, nil)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)
	return result, err
}
