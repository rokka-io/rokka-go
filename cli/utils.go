package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func PrettyJSON(data interface{}) string {
	var s []byte
	var err error

	switch data := data.(type) {
	case []byte:
		var pretty bytes.Buffer
		err = json.Indent(&pretty, data, "", "  ")
		s = pretty.Bytes()
	default:
		s, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling JSON: %s\n", err)
		os.Exit(1)
	}

	return string(s)
}
