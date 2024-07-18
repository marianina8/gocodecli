package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	startDate string
)

type LogEntry struct {
	Time       time.Time `json:"time"`
	URL        string    `json:"url"`
	StatusCode int       `json:"statusCode"`
	Duration   int64     `json:"duration"`
}

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history [urls]",
	Short: "Displays the history of health checks for specified URL(s)",
	Long: `The history command parses a log file for historical data related to specific URL checks.
The --startDate flag can be used to specify the UTC start date of the history period.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		_, err := time.Parse("01/02/2006", startDate)
		if err != nil {
			fmt.Printf("Error parsing startDate flag value: %v\n", err)
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		displayHistory(args)
	},
}

func init() {
	historyCmd.Flags().StringVar(&startDate, "startDate", "", "The start date for displaying history (format: MM/DD/YYYY)")
	rootCmd.AddCommand(historyCmd)
}

func displayHistory(urls []string) {
	urlMap := toURLMap(urls)
	file, err := os.Open(logFile)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	startDateParsed, err := time.Parse("01/02/2006", startDate)
	if err != nil {
		fmt.Printf("Error parsing startDate: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			fmt.Printf("Error unmarshalling log line: %v\n", err)
			continue
		}

		if urlMap[entry.URL] && entry.Time.After(startDateParsed) {
			fmt.Println(line)
		}
	}
}

func toURLMap(urls []string) map[string]bool {
	urlMap := make(map[string]bool)
	for _, url := range urls {
		urlMap[url] = true
	}
	return urlMap
}
