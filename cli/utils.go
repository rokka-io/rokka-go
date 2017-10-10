package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrettyPrintJSON(data interface{}) {
	s, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(s))
}
