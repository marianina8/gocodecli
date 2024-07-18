package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var ExitFunction = os.Exit
var outputWriter io.Writer = os.Stdout

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the health of specified URL(s)",
	Long:  `Performs a health check by sending a request to the specified URL(s) and reports the status.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if output == "table" {
			// print table
			table := tablewriter.NewWriter(outputWriter)
			upColor := color.New(color.FgGreen).SprintFunc()
			downColor := color.New(color.FgRed).SprintFunc()
			table.SetHeader([]string{"URL", "Status"})
			for _, url := range args {
				status := downColor("Down")
				if checkURL(ctx, url, threshold, retries) {
					status = upColor("Up")
				}
				table.Append([]string{url, status})
			}
			table.Render()
		} else {
			for _, url := range args {
				checkURL(ctx, url, threshold, retries)
			}
		}
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
	rootCmd.AddCommand(checkCmd)
}

func checkURL(ctx context.Context, url string, threshold float64, retries int) bool {
	var resp *http.Response
	var lastError error
	var duration time.Duration

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	if output != "table" || silent {
		s.Start()
	}
	for attempt := 0; attempt <= retries; attempt++ {
		start := time.Now()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			l.ErrorContext(ctx, "failed to create request", "url", url)
		}

		client := &http.Client{}
		resp, err = client.Do(req)
		s.Disable()
		if err != nil {
			l.ErrorContext(ctx, "failed to perform request", "url", url, "err", err)
		}

		duration = time.Since(start)
		if err == nil {
			defer resp.Body.Close()
			if duration.Seconds() > threshold {
				l.WarnContext(ctx, "exceeded threshold", "url", url, "responseTime", duration)
			} else {
				l.InfoContext(ctx, "successful check", "url", url, "statusCode", resp.StatusCode, "duration", duration)
			}
			break
		}

		select {
		case <-ctx.Done():
			l.ErrorContext(ctx, "check cancelled", "url", url, "attempt", attempt, "err", err)
			ExitFunction(1)
		case <-time.After(2 * time.Second):
			l.InfoContext(ctx, "backing off", "url", url)
		}
		lastError = err
	}
	if lastError != nil {
		l.ErrorContext(ctx, "fetching error", "url", url, "retries", retries, "err", lastError)
		return false
	}

	if resp.StatusCode == http.StatusOK {
		return true
	} else {
		return false
	}

}

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}
