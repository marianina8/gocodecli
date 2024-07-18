package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type URLList []string

func (list *URLList) String() string {
	return fmt.Sprint(*list)
}

func (list *URLList) Set(value string) error {
	for _, url := range strings.Split(value, ",") {
		*list = append(*list, url)
	}
	return nil
}

var (
	url       string
	logFile   string
	logger    *log.Logger
	interval  time.Duration
	silent    bool
	verbose   bool
	threshold float64
	retries   int
	urls      URLList
)

func main() {
	flag.StringVar(&url, "url", "", "URL to check")
	flag.StringVar(&logFile, "logfile", "healthcheck.log", "File to log output to")
	flag.DurationVar(&interval, "interval", 2*time.Second, "Interval between healthchecks")
	flag.BoolVar(&silent, "silent", false, "Run in silent mode without stdout output")
	flag.BoolVar(&verbose, "verbose", false, "Run in verbose mode.  Overrides silent mode")
	flag.Float64Var(&threshold, "threshold", 0.5, "Threshold value for considering a response to be too slow (in seconds)")
	flag.IntVar(&retries, "retries", 3, "Number of retries for a failed request")
	flag.Var(&urls, "urls", "Comma-separated list of URLs to check")
	flag.Usage = customUsage
	flag.Parse()

	if len(os.Args) <2 {
		fmt.Println ("No flags provided.")
		flag.Usage()
		os.Exit(1)
	}
	if (url != "" && len(urls) > 0) || (url == "" && len(urls) == 0) {
		log.Fatalf("Specify either a single URL using -url or multiple URLS using -urls, but not both.")
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logfile %s: %v", logFile, err)
		log.Fatal(err)
	}
	defer file.Close()

	if silent && !verbose {
		logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logger = log.New(io.MultiWriter(file, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	ticker := time.NewTicker(interval)
	for range ticker.C {
		if url != "" {
			checkURL(url, threshold, retries)
		} else {
			for _, u := range urls {
				checkURL(u, threshold, retries)
			}
		}
	}

}

func checkURL(url string, threshold float64, retries int) {
	var resp *http.Response
	var err error
	var duration time.Duration
	for attempt := 0; attempt <= retries; attempt++ {
		start := time.Now()
		resp, err = http.Get(url)
		duration = time.Since(start)
		if err == nil {
			defer resp.Body.Close()
			break
		}

		if attempt < retries {
			fmt.Fprintf(os.Stderr, "Attempt %d failed, retrying...\n", attempt+1)
			time.Sleep(time.Second * 2) // Backoff strategy
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching %s after %d retries: %v\n", url, retries, err)
		return
	}

	if duration.Seconds() > threshold && verbose {
		fmt.Fprintf(os.Stderr, "Warning: %s response time (%v) exceeded threshold of %fs\n", url, duration, threshold)
	}
	logger.Printf("Checked %s, Status: %d, Response Time: %v\n", url, resp.StatusCode, duration)
}

func customUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "This tool performs health checks on specified URLS.\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Options: \n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nExamples:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  %s -url https://example.com\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "  %s -urls https://example.com,https://example2.com\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "  %s -url https://example.com -logfile results.log\n", os.Args[0])
}