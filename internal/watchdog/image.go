package watchdog

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/davidegagliardi/syowatchdog/pkg/types"
)

type ImageProcessor struct {
	httpClient *http.Client
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (ip *ImageProcessor) FetchAndProcessImage(url string) (*types.ImageData, error) {
	// Fetch image from URL
	imageURL := strings.Trim(url, "\"")
	resp, err := ip.httpClient.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Convert to base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)

	// Calculate hash for comparison
	hash := fmt.Sprintf("%x", md5.Sum(imageData))

	return &types.ImageData{
		Base64Content: base64Data,
		Hash:          hash,
		LastUpdated:   time.Now(),
		SourceURL:     imageURL,
	}, nil
}

func (ip *ImageProcessor) CompareImages(oldData, newData *types.ImageData) bool {
	if oldData == nil {
		return true // First time, consider it changed
	}
	return oldData.Hash != newData.Hash
}
