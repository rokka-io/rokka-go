package massupload

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"
)

type ResultReporter struct {
	template *template.Template
	writer   *tabwriter.Writer
}

func BuildReporter(tpl string, out *os.File) (ResultReporter, error) {

	reporter := ResultReporter{}
	templ, err := template.New("").Parse(tpl)
	if err != nil {
		return reporter, fmt.Errorf("error parsing response template: %s", err)
	}

	reporter.template = templ
	reporter.writer = tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)

	return reporter, nil
}

func (r *ResultReporter) Report(item interface{}) error {
	err := r.template.Execute(r.writer, item)
	if err != nil {
		return err
	}
	r.writer.Flush()

	return nil
}
