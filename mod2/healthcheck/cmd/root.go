package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/marianina8/gocodecli/mod2/healthcheck/logger"
	"github.com/spf13/cobra"
)

var (
	logFile string
	l       *slog.Logger

	threshold float64
	retries   int

	silent      bool
	verbose     bool
	versionFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "A tool for monitoring health status and responsiveness of web applications",
	Long: `The healthcheck command is designed to assess the health and
responsiveness of specified web applications. It sends HTTP requests
to URLs provided by the user, evaluating whether the services are
accessible and how quickly they respond. This command supports both
immediate, one-off checks and continuous monitoring, allowing users
to specify intervals for ongoing health assessments. With additional
flags for customization, users can tailor the command to meet various
monitoring needs, from simple uptime checks to detailed performance
analysis."`,
	Run: func(cmd *cobra.Command, args []string) {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			printVersion()
			os.Exit(0)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		silent, _ := cmd.Flags().GetBool("silent")
		l = logger.New(logFile, verbose, silent)
	},
}

func Execute() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logFile, "logfile", "healthcheck.log", "File to log output to")
	rootCmd.PersistentFlags().Float64Var(&threshold, "threshold", 0.5, "Threshold value for considering a response to be too slow (in seconds)")
	rootCmd.PersistentFlags().IntVar(&retries, "retries", 3, "Number of retries for a failed request")
	rootCmd.PersistentFlags().BoolVar(&silent, "silent", false, "Run in silent mode without stdout output")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Run in verbose mode.  Overrides silent mode")
	rootCmd.Flags().BoolVar(&versionFlag, "version", false, "Print version")

}
