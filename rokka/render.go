package rokka

import (
	"fmt"
	"strings"
)

func (c *Client) GetURL(organization, hash, format string, ops []Operation) string {
	host := strings.Replace(c.config.ImageHost, "{{organization}}", organization, -1)

	if len(ops) == 0 {
		ops = append(ops, NoopOperation{})
	}
	opURL := make([]string, len(ops))
	for i, o := range ops {
		opURL[i] = o.toURLPath()
	}

	return fmt.Sprintf("%s/dynamic/%s.%s", host, strings.Join(opURL, "--"), format)
}
