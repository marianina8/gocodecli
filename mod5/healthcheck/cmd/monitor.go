package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	interval time.Duration
)

var monitorCmd = &cobra.Command{
	Use:   "monitor [urls]",
	Short: "Monitor the health of specified URL(s) over time",
	Long:  `Continuously monitors the health of the specified URL(s) at the specified interval`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		monitorURLs(ctx, args)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		for _, url := range args {
			if !isValidURL(url){
				return fmt.Errorf("The URL %s is invalid.  Please ensure it includes the protocol (http or https) and a valid domain.", url)
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
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	if output == "table" {
		s.Start()
	}
	for range ticker.C {
		if output == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"URL", "Status", "Last Time Checked"})

			upColor := color.New(color.FgGreen).SprintFunc()
			downColor := color.New(color.FgRed).SprintFunc()
			for _, url := range urls {
				status := downColor("Down")
				if checkURL(ctx, url, threshold, retries) {
					status = upColor("Up")
				}
				currentTime := time.Now()
				formattedTime := currentTime.Format("01/02/2006 03:04PM")
				table.Append([]string{url, status, formattedTime})
			}
			// Clear the screen
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				fmt.Println("Unable to clear the screen: ", err)
			}
			s.Disable()
			table.Render()
		} else {
			for _, url := range urls {
				checkURL(ctx, url, threshold, retries)
			}
		}
	}
}
