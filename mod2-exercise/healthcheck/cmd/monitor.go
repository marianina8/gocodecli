package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	interval time.Duration
)

var monitorCmd = &cobra.Command{
	Use:   "monitor [urls]",
	Short: "Monitor the health of specified URL(s) over time",
	Long:  `Continuously monitors the health of the specified URL(s) at the specified interval`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		monitorURLs(ctx, args)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		for _, url := range args {
			if !isValidURL(url) {
				return fmt.Errorf("invalid url: %s", url)
			}
		}
		return nil
	},
}

func init() {
	monitorCmd.Flags().DurationVar(&interval, "interval", 2*time.Second, "Interval between healthchecks")
	rootCmd.AddCommand(monitorCmd)
}

func monitorURLs(ctx context.Context, urls []string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for _, url := range urls {
			checkURL(ctx, url, threshold, retries)
		}
	}
}
