package rokka

import (
	"fmt"
	"strings"
)

// GetURL generates an URL for the given organization, image hash, and format based on the list of operations
// given to it.
// If the operation slice is empty it generates a default noop operation.
// The URL returned has the format: `https://{imageHost for org}/dynamic/{stacks-applied-with-options}/{hash}.{format}`.
func (c *Client) GetURL(organization, hash, format string, ops []Operation) (string, error) {
	return c.GetURLForStack(organization, hash, format, "dynamic", ops)
}

// GetURLForStack generates an URL for the given organization, stack name, image hash, and format.
// The operations are added on top of the stack.
// If the operation slice is empty it generates a default noop operation.
// The URL returned has the format: `https://{imageHost for org}/{stack}/{stacks-applied-with-options}/{hash}.{format}`.
func (c *Client) GetURLForStack(organization, hash, format, stack string, ops []Operation) (string, error) {
	host := strings.Replace(c.config.ImageHost, "{{organization}}", organization, -1)

	if len(ops) == 0 {
		ops = append(ops, &NoopOperation{})
	}
	opURL := make([]string, len(ops))
	for i, o := range ops {
		if _, err := o.Validate(); err != nil {
			return "", err
		}

		opURL[i] = o.toURLPath()
	}

	return fmt.Sprintf("%s/%s/%s/%s.%s", host, stack, strings.Join(opURL, "--"), hash, format), nil
}
