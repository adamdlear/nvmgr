package metadata

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type ConfigMetadata struct {
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description,omitempty"`
}

func MetadataPath(configDir string) string {
	return filepath.Join(configDir, ".nvmgr.json")
}

func Write(configDir, name, desc string) error {
	meta := ConfigMetadata{
		Name:        name,
		CreatedAt:   time.Now(),
		Description: desc,
	}
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(MetadataPath(configDir), data, 0o644)
}

func Read(configDir string) (*ConfigMetadata, error) {
	data, err := os.ReadFile(MetadataPath(configDir))
	if err != nil {
		return nil, err
	}
	var meta ConfigMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}
