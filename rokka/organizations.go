package rokka

type OrganizationResponse struct{}

func (c *Client) GetOrganization(name string) (OrganizationResponse, error) {
	result := OrganizationResponse{}

	req, err := c.NewRequest("GET", "/organizations/"+name, nil, nil)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)
	return result, err
}
