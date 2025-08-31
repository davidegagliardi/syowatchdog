package types

import "time"

type ImageData struct {
	Base64Content string    `json:"base64_content"`
	Hash          string    `json:"hash"`
	LastUpdated   time.Time `json:"last_updated"`
	SourceURL     string    `json:"source_url"`
}
