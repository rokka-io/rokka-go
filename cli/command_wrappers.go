package cli

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/rokka-io/rokka-go/rokka"
)

func setDefaultValue(data map[string]string, key, val string) {
	if _, ok := data[key]; !ok {
		data[key] = val
	}
}

func login(c *rokka.Client, _ map[string]string, options map[string]string) (interface{}, error) {
	valid, err := c.ValidAPIKey()
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid API key")
	}

	cfg := Config{
		APIKey: c.GetConfig().APIKey,
	}
	err = SaveConfig(cfg)
	if err != nil {
		return nil, err
	}
	return "Login successful", nil
}

func getStackOptions(c *rokka.Client, args, options map[string]string) (interface{}, error) {
	return c.GetStackOptions()
}

func getOrganization(c *rokka.Client, args, options map[string]string) (interface{}, error) {
	return c.GetOrganization(args["org"])
}

func createOrganization(c *rokka.Client, args, options map[string]string) (interface{}, error) {
	return c.CreateOrganization(args["org"], args["billingEmail"], args["displayName"])
}

func listSourceImages(c *rokka.Client, args, options map[string]string) (interface{}, error) {
	setDefaultValue(options, "limit", "20")

	return c.ListSourceImages(args["org"], options)
}

func getSourceImage(c *rokka.Client, args, options map[string]string) (interface{}, error) {
	return c.GetSourceImage(args["org"], args["hash"])
}

func createSourceImage(c *rokka.Client, args map[string]string, options map[string]string) (interface{}, error) {
	if _, err := os.Stat(args["file"]); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(args["file"])
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return c.CreateSourceImage(args["org"], path.Base(args["file"]), file)
}

func listStacks(c *rokka.Client, args, options map[string]string) (interface{}, error) {
	return c.ListStacks(args["org"])
}

type dailyStats struct {
	Date                     time.Time
	Space, Files, Downloaded int
}

func getStats(c *rokka.Client, args, options map[string]string) (interface{}, error) {
	setDefaultValue(options, "from", time.Now().Add(-30*24*time.Hour).Format("2006-01-02"))
	setDefaultValue(options, "to", time.Now().Format("2006-01-02"))

	from, err := time.Parse("2006-01-02", options["from"])
	if err != nil {
		return nil, fmt.Errorf(`Invalid format for parameter "from". Expected YYYY-MM-DD, got "%s"`, options["from"])
	}
	to, err := time.Parse("2006-01-02", options["to"])
	if err != nil {
		return nil, fmt.Errorf(`Invalid format for parameter "to". Expected YYYY-MM-DD, got "%s"`, options["to"])
	}

	res, err := c.GetStats(args["org"], options)
	if err != nil {
		return nil, err
	}

	result := make(map[string]dailyStats)
	for i, max := 0, int(to.Sub(from).Hours()/24); i <= max; i++ {
		d := from.Add(time.Duration(i) * 24 * time.Hour)
		result[d.Format("2006-01-02")] = dailyStats{Date: d}
	}
	for _, v := range res.SpaceInBytes {
		d := v.Timestamp.Format("2006-01-02")
		r := result[d]
		r.Space = v.Value
		result[d] = r
	}
	for _, v := range res.NumberOfFiles {
		d := v.Timestamp.Format("2006-01-02")
		r := result[d]
		r.Files = v.Value
		result[d] = r
	}
	for _, v := range res.BytesDownloaded {
		d := v.Timestamp.Format("2006-01-02")
		r := result[d]
		r.Downloaded = v.Value
		result[d] = r
	}

	return result, nil
}
