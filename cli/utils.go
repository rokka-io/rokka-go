package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// PrettyJSON marshals JSON with an indent.
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

// prettyDateTime formats a time.Time.
func prettyDateTime(t time.Time) string {
	return t.Format("Mon, January _2 2006, 15:04:05Z0700")
}

// prettyDate formats a time.Time.
func prettyDate(t time.Time) string {
	return t.Format("Jan _2 2006")
}

func join(d interface{}, sep string) (string, error) {
	switch d := d.(type) {
	case map[string]interface{}:
		r := make([]string, 0)
		for k := range d {
			r = append(r, k)
		}
		return strings.Join(r, sep), nil
	}
	return "", errors.New("cli/join: unhandled type to join with separator")
}
