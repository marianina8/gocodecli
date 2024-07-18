package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the health of specified URL(s)",
	Long:  `Performs a health check by sending a request to the specified URL(s) and reports the status.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		for _, url := range args {
			checkURL(ctx, url, threshold, retries)
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

func checkURL(ctx context.Context, url string, threshold float64, retries int) {
	var resp *http.Response
	var lastError error
	var duration time.Duration
	for attempt := 0; attempt <= retries; attempt++ {
		start := time.Now()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			l.ErrorContext(ctx, "failed to create request", "url", url)
		}

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			l.ErrorContext(ctx, "failed to perform request", "url", url, "err", err)
		}

		duration = time.Since(start)
		if err == nil {
			defer resp.Body.Close()
			if duration.Seconds() >threshold {
				l.WarnContext(ctx, "exceeded threshold", "url", url, "responseTime", duration)
			} else {
				l.InfoContext(ctx, "successful check", "url", url, "statusCode", resp.StatusCode, "duration", duration)
			}
			break
		}

		select {
		case <-ctx.Done():
			l.ErrorContext(ctx, "check cancelled", "url", url, "attempt", attempt, "err", err)
			os.Exit(1)
		case <-time.After(2 * time.Second):
			l.InfoContext(ctx, "backing off", "url", url)
		}

		lastError = err
	}
	if lastError != nil {
		l.ErrorContext(ctx, "fetching error", "url", url, "retries", retries, "err", lastError)
		return
	}

}

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}
