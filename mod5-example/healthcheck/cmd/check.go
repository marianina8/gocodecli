package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var ExitFunction = os.Exit
var outputWriter io.Writer = os.Stdout

type URLValidationError struct {
	URL    string
	Detail string
}

func (e *URLValidationError) Error() string {
	return fmt.Sprintf("The URL %s is invalid. Details: %s", e.URL, e.Detail)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the health of specified URL(s)",
	Long:  `Performs a health check by sending a request to the specified URL(s) and reports the status.`,
	//Args:  cobra.MinimumNArgs(1),
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
		if len(args) == 0 {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Enter URLs to check, one per line.  Press Enter twice to finish:")
			for {
				fmt.Print("Enter URL:")
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("Error parsing input string, details: %v", err)
				}
				input = strings.TrimSpace(input)
				if input == "" {
					break
				}
				err = isValidURL(input)
				if err != nil {
					return err
				}
			}
		}
		for _, url := range args {
			err := isValidURL(url)
			if err != nil {
				return err
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
			return false
		}

		l.DebugContext(ctx, "request details", "method", req.Method, "url", req.URL, "headers", req.Header)

		client := &http.Client{}
		resp, err = client.Do(req)
		s.Disable()
		if err != nil {
			l.ErrorContext(ctx, "failed to perform request", "url", url, "err", err)
			return false
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

func isValidURL(u string) error {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return &URLValidationError{
			URL:    u,
			Detail: err.Error(),
		}
	}
	if parsedURL.Scheme == "" {
		return &URLValidationError{
			URL:    u,
			Detail: "missing scheme",
		}
	}
	if parsedURL.Host == "" {
		return &URLValidationError{
			URL:    u,
			Detail: "missing host",
		}
	}
	return nil
}
