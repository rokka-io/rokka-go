package cli

import (
	"fmt"
	"time"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

type dailyStats struct {
	Date       time.Time
	Space      int
	Files      int
	Downloaded int
}

var statsGetCmdOptions rokka.GetStatsOptions

func getStats(c *rokka.Client, args []string) (interface{}, error) {
	from, err := time.Parse("2006-01-02", statsGetCmdOptions.From)
	if err != nil {
		return nil, fmt.Errorf(`Invalid format for parameter "from". Expected YYYY-MM-DD, got "%s"`, statsGetCmdOptions.From)
	}
	to, err := time.Parse("2006-01-02", statsGetCmdOptions.To)
	if err != nil {
		return nil, fmt.Errorf(`Invalid format for parameter "to". Expected YYYY-MM-DD, got "%s"`, statsGetCmdOptions.To)
	}
	res, err := c.GetStats(args[0], statsGetCmdOptions)
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

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:                   "stats",
	Short:                 "Gather statistics",
	Run:                   nil,
	Aliases:               []string{"st"},
	DisableFlagsInUseLine: true,
}

var statsGetCmd = &cobra.Command{
	Use:                   "get [org]",
	Short:                 "Show statistics for an organization (traffic, space used, files uploaded)",
	Args:                  cobra.ExactArgs(1),
	Aliases:               []string{"g"},
	DisableFlagsInUseLine: true,
	Run:                   run(getStats, "Date\tDownloaded (Bytes)\tSpace (Bytes)\tFiles\n{{range $_, $e := .}}{{date $e.Date}}\t{{ $e.Downloaded }}\t{{ $e.Space }}\t{{ $e.Files }}\n{{end}}"),
}

func init() {
	rootCmd.AddCommand(statsCmd)

	statsCmd.AddCommand(statsGetCmd)

	statsGetCmd.Flags().StringVar(&statsGetCmdOptions.From, "from", time.Now().Add(-30*24*time.Hour).Format("2006-01-02"), "Start date, format: YYYY-MM-DD")
	statsGetCmd.Flags().StringVar(&statsGetCmdOptions.To, "to", time.Now().Format("2006-01-02"), "End date, format: YYYY-MM-DD")
}
