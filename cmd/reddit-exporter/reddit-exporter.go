package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/travisgroth/reddit-exporter/internal/scanner"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

var verbose bool
var subs []string
var port int
var address string
var regexFiles []string

var rootCmd = &cobra.Command{
	Use:   "reddit-exporter",
	Short: "Exports prometheus style metrics about reddit",
	Long:  "Configurable prometheus metric exporter for reddit and subreddit activity",
	Run: func(cmd *cobra.Command, args []string) {

		if verbose {
			log.SetLevel(log.DebugLevel)
		}

		cfg := graw.Config{SubredditComments: subs, Subreddits: subs}
		script, _ := reddit.NewScript("graw:reddit-exporter:0.1.0", time.Second*1)

		scanner := scanner.NewScanner(
			cfg,
			script,
			graw.Scan,
		)

		for _, regexFile := range regexFiles {
			err := addRegexHandlers(scanner, regexFile)
			if err != nil {
				log.Fatalf("Failed to load regex file %s: %s", regexFile, err)
			}
		}

		http.Handle("/metrics", promhttp.Handler())
		bindAddress := fmt.Sprintf("%s:%d", address, port)
		log.Info("Listening on ", bindAddress)
		go func() { log.Fatal(http.ListenAndServe(bindAddress, nil)) }()

		scanner.Run()

	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Turn on verbose mode")
	rootCmd.PersistentFlags().StringSliceVarP(&subs, "subreddit", "s", nil,
		"Subreddit(s) to monitor (Required).\nCan be specified multiple times or comma separated")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8000, "Metrics server port")
	rootCmd.PersistentFlags().StringVar(&address, "address", "0.0.0.0", "Metrics server bind address")
	rootCmd.PersistentFlags().StringSliceVar(&regexFiles, "regexfile", nil,
		"File containing regex matches in format 'name;regex'\nCan be specified multiple times or comma separated")

	rootCmd.MarkPersistentFlagRequired("subreddit")
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	rootCmd.Execute()
}
