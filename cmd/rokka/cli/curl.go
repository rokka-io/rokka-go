package cli

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

var curlOptions struct {
	method         string
	headers        []string
	data           string
	includeHeaders bool
}

type curlData struct {
	Headers    http.Header
	StatusLine string
	Data       string
}

func curlResponseHandler(resp *http.Response, v interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cd := v.(*curlData)
	cd.Data = string(body)
	cd.Headers = resp.Header
	cd.StatusLine = fmt.Sprintf("%s %s", resp.Proto, resp.Status)

	return nil
}

func curl(c *rokka.Client, args []string) (interface{}, error) {
	var body io.Reader
	var err error
	if curlOptions.data != "" {
		if strings.HasPrefix(curlOptions.data, "@") {
			body, err = os.Open(curlOptions.data[1:])
			if err != nil {
				return nil, err
			}
		} else {
			body = bytes.NewBufferString(curlOptions.data)
		}
	}

	req, err := c.NewRequest(curlOptions.method, args[0], body, nil)
	if err != nil {
		return nil, err
	}

	for _, hdr := range curlOptions.headers {
		parts := strings.Split(hdr, ":")
		if len(parts) == 2 {
			req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	res := curlData{}

	if err := c.Call(req, &res, curlResponseHandler); err != nil {
		return nil, err
	}

	if !curlOptions.includeHeaders {
		res.Headers = nil
		res.StatusLine = ""
	}

	return res, nil
}

// curlCmd allows to use a curl-like interface for simple requests
var curlCmd = &cobra.Command{
	Use:                   "curl [path]",
	Short:                 "cURL-like interface for doing raw requests.",
	Run:                   run(curl, "{{if .StatusLine}}{{.StatusLine}}\n{{end}}{{if .Headers}}{{range $name, $val := .Headers}}{{range $val}}{{$name}}: {{.}}\n{{end}}{{end}}\n{{end}}{{.Data}}"),
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(curlCmd)

	curlCmd.Flags().StringVarP(&curlOptions.method, "request", "X", "GET", "HTTP method (GET,POST,PUT, etc.)")
	curlCmd.Flags().StringArrayVarP(&curlOptions.headers, "header", "H", nil, "Custom headers to add to the request")
	curlCmd.Flags().StringVarP(&curlOptions.data, "data", "d", "", "HTTP data to include in request body (or a filename prefixed with '@')")
	curlCmd.Flags().BoolVarP(&curlOptions.includeHeaders, "include", "i", false, "Include HTTP response headers in the output")
}
