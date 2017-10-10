package main

import (
	"log"
	"encoding/json"
	"github.com/rokka-io/rokka-go/client"
)

func main() {
	c := client.NewClient(&client.Config{APIKey: ""})
	res, err := c.GetStackoptions()
	
	if err != nil {
		log.Fatal(err)
	}

	pretty, _ := json.MarshalIndent(res, "", "  ")
	println(string(pretty))
}
