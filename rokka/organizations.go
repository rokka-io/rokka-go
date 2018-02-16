package rokka

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// OrganizationResponse contains the information about an organization.
type OrganizationResponse struct {
	ID                 string    `json:"id"`
	DisplayName        string    `json:"display_name"`
	Name               string    `json:"name"`
	BillingEmail       string    `json:"billing_email"`
	Created            time.Time `json:"created"`
	MasterOrganization string    `json:"master_organization"`
	Limit              *struct {
		SpaceInBytes   int `json:"space_in_bytes,omitempty"`
		TrafficInBytes int `json:"traffic_in_bytes,omitempty"`
	} `json:"limit,omitempty"`
}

type createOrganizationRequest struct {
	DisplayName  string `json:"display_name,omitempty"`
	BillingEmail string `json:"billing_email"`
}

// GetOrganization returns information about the requested organization if the passed API key is allowed to access it.
//
// See: https://rokka.io/documentation/references/organizations.html#read-data-of-one-organization
func (c *Client) GetOrganization(name string) (OrganizationResponse, error) {
	result := OrganizationResponse{}

	req, err := c.NewRequest(http.MethodGet, "/organizations/"+name, nil, nil)
	if err != nil {
		return result, err
	}

	err = c.CallJSONResponse(req, &result)
	return result, err
}

// CreateOrganization can be used to create a new organization which is bound to the API key (user) you supplied.
//
// See: https://rokka.io/documentation/references/organizations.html#create-an-organization
func (c *Client) CreateOrganization(name, billingEmail, displayName string) (OrganizationResponse, error) {
	result := OrganizationResponse{}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(createOrganizationRequest{
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

	err = c.CallJSONResponse(req, &result)
	return result, err
}
