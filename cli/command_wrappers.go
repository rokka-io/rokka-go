package cli

import (
	"fmt"

	"github.com/rokka-io/rokka-go/rokka"
)

func getStackoptions(c *rokka.Client, _ map[string]string, _ map[string]string) error {
	res, err := c.GetStackoptions()
	if err != nil {
		return err
	}
	fmt.Println(PrettyJSON(res))

	return nil
}

func getOrganization(c *rokka.Client, args map[string]string, _ map[string]string) error {
	res, err := c.GetOrganization(args["name"])
	if err != nil {
		return err
	}
	fmt.Println(PrettyJSON(res))

	return nil
}
