package cli

import (
	"fmt"

	"github.com/rokka-io/rokka-go/rokka"
)

func getStackOptions(c *rokka.Client, logger *Log, _ map[string]string, _ map[string]string) error {
	res, err := c.GetStackOptions()
	if err != nil {
		return err
	}
	s, err := PrettyJSON(res)
	if err != nil {
		return fmt.Errorf("cli/getStackOptions: Error pretty printing JSON: %s", err)
	}
	logger.Printf("%s", s)

	return nil
}

func getOrganization(c *rokka.Client, logger *Log, args map[string]string, _ map[string]string) error {
	res, err := c.GetOrganization(args["name"])
	if err != nil {
		return err
	}
	s, err := PrettyJSON(res)
	if err != nil {
		return fmt.Errorf("cli/getOrganization: Error pretty printing JSON: %s", err)
	}
	logger.Printf("%s", s)

	return nil
}
