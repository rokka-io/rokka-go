package rokka

import (
	"fmt"
	"strings"
)

func (c *Client) GetURL(organization, hash, format string, ops []Operation) (string, error) {
	host := strings.Replace(c.config.ImageHost, "{{organization}}", organization, -1)

	if len(ops) == 0 {
		ops = append(ops, NoopOperation{})
	}
	opURL := make([]string, len(ops))
	for i, o := range ops {
		if _, err := o.Validate(); err != nil {
			return "", err
		}

		opURL[i] = o.toURLPath()
	}

	return fmt.Sprintf("%s/dynamic/%s/%s.%s", host, strings.Join(opURL, "--"), hash, format), nil
}
