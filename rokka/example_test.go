package rokka

import (
	"fmt"
	"os"
)

func Example() {
	c := NewClient(&Config{
		APIKey: "exampleAPIKey",
	})

	resp, err := c.GetOrganization("example")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s", resp.ID)
}
