package watchdog

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davidegagliardi/syowatchdog/internal/config"
	"github.com/davidegagliardi/syowatchdog/internal/storage"
	"github.com/davidegagliardi/syowatchdog/internal/telegram"
)

type Watchdog struct {
	config         *config.Config
	imageProcessor *ImageProcessor
	telegramClient *telegram.Client
	storage        *storage.FileStorage
	ticker         *time.Ticker
	stopChan       chan struct{}
}

func New(cfg *config.Config) *Watchdog {
	return &Watchdog{
		config:         cfg,
		imageProcessor: NewImageProcessor(),
		telegramClient: telegram.NewClient(cfg.TelegramBotToken, cfg.TelegramChatID),
		storage:        storage.NewFileStorage(cfg.StoragePath),
		stopChan:       make(chan struct{}),
	}
}

func (w *Watchdog) Start() error {
	log.Printf("Starting Syowatchdog...")
	log.Printf("Monitoring image: %s", w.config.ImageURL)
	log.Printf("Check interval: %v", w.config.CheckInterval)

	// Test Telegram connection
	if err := w.telegramClient.TestConnection(); err != nil {
		return fmt.Errorf("failed to connect to Telegram: %w", err)
	}
	log.Println("Telegram connection verified")

	// Send startup notification
	if err := w.telegramClient.SendMessage("🤖 Syowatchdog started monitoring your image!"); err != nil {
		log.Printf("Warning: Failed to send startup notification: %v", err)
	}

	// Initial check
	if err := w.performImageCheck(); err != nil {
		log.Printf("Initial check failed: %v", err)
	}

	// Start periodic checks
	w.ticker = time.NewTicker(w.config.CheckInterval)
	defer w.ticker.Stop()

	for {
		select {
		case <-w.ticker.C:
			if err := w.performImageCheck(); err != nil {
				log.Printf("Image check failed: %v", err)
			}
		case <-w.stopChan:
			log.Println("Watchdog stopped")
			return nil
		}
	}
}

func (w *Watchdog) StartWithGracefulShutdown(ctx context.Context) error {
	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start watchdog in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- w.Start()
	}()

	// Wait for shutdown signal or error
	select {
	case <-sigChan:
		log.Println("Received shutdown signal")
		return w.Stop()
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return w.Stop()
	}
}

func (w *Watchdog) Stop() error {
	log.Println("Stopping Syowatchdog...")

	// Send shutdown notification
	if err := w.telegramClient.SendMessage("🛑 Syowatchdog stopped monitoring."); err != nil {
		log.Printf("Warning: Failed to send shutdown notification: %v", err)
	}

	close(w.stopChan)
	if w.ticker != nil {
		w.ticker.Stop()
	}
	return nil
}

func (w *Watchdog) performImageCheck() error {
	log.Println("Checking image for changes...")

	// Fetch and process new image
	newImageData, err := w.imageProcessor.FetchAndProcessImage(w.config.ImageURL)
	if err != nil {
		return fmt.Errorf("failed to process image: %w", err)
	}

	// Load previous image data
	oldImageData, err := w.storage.LoadImageData()
	if err != nil {
		return fmt.Errorf("failed to load previous image data: %w", err)
	}

	// Compare images
	if w.imageProcessor.CompareImages(oldImageData, newImageData) {
		log.Println("Image change detected!")

		// Send notification
		if err := w.telegramClient.SendImageChangeNotification(
			w.config.ImageURL,
			newImageData.LastUpdated,
		); err != nil {
			log.Printf("Failed to send Telegram notification: %v", err)
		}

		// Save new image data
		if err := w.storage.SaveImageData(newImageData); err != nil {
			return fmt.Errorf("failed to save image data: %w", err)
		}

		log.Println("Image data updated and notification sent")
	} else {
		log.Println("No changes detected")

		// Update the last checked time even if no changes
		if err := w.storage.SaveImageData(newImageData); err != nil {
			log.Printf("Warning: Failed to update timestamp: %v", err)
		}
	}

	return nil
}
