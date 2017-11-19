package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

func PrettyJSON(data interface{}) (string, error) {
	var err error
	var pretty bytes.Buffer

	switch data := data.(type) {
	case []byte:
		err = json.Indent(&pretty, data, "", "  ")
	default:
		enc := json.NewEncoder(&pretty)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		err = enc.Encode(data)
	}

	if err != nil {
		return "", fmt.Errorf("cli/prettyJSON: %s", err)
	}

	return pretty.String(), nil
}

func PrettyDate(t time.Time) string {
	return t.Format("Mon, January _2 2006, 15:04:05Z0700")
}
