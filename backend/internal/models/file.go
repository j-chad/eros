package models

import "time"

type File struct {
	ID     string `json:"id"`
	NodeID string `json:"node_id"`

	Filename   string `json:"filename"`
	MimeType   string `json:"mime_type"`
	StorageKey string `json:"storage_key"`
	SizeBytes  int64  `json:"size_bytes"`

	CreatedAt time.Time `json:"created_at"`
}
