package cmd

import (
	"context"
	"log"

	"github.com/davidegagliardi/syowatchdog/internal/config"
	"github.com/davidegagliardi/syowatchdog/internal/watchdog"
	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "syowatchdog",
	Short: "A watchdog CLI tool for monitoring image changes",
	Long: `Syowatchdog monitors an image URL for changes and sends notifications via Telegram.
	
The tool fetches an image, converts it to base64, stores it persistently, and checks
for changes at regular intervals. When a change is detected, it sends a notification
via Telegram.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		if err := cfg.Validate(); err != nil {
			log.Fatalf("Invalid configuration: %v", err)
		}

		wd := watchdog.New(cfg)
		if err := wd.StartWithGracefulShutdown(context.Background()); err != nil {
			log.Fatalf("Watchdog failed: %v", err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
