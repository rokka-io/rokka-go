package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrettyJSON(data interface{}) (string, error) {
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
		return "", fmt.Errorf("cli/prettyJSON: %s", err)
	}

	return string(s), nil
}
