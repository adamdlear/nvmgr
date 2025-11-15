package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Config represents a single nvim configuration managed by nvmgr.
type Config struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	// ConfigPrefix is the prefix for all nvmgr managed configurations.
	ConfigPrefix = "nvim-"
)

// ConfigDir returns the directory where nvmgr stores its configurations.
func ConfigDir() string {
	return filepath.Join(os.Getenv("HOME"), ".config")
}

// ConfigPath returns the full path to a nvmgr managed configuration.
func ConfigPath(name string) string {
	return filepath.Join(ConfigDir(), ConfigPrefix+name)
}

func nvmgrConfigsPath() string {
	return filepath.Join(ConfigDir(), "nvmgr", "configs.json")
}

// Exists returns true if a configuration with the given name exists.
func Exists(name string) bool {
	_, err := os.Stat(ConfigPath(name))
	return err == nil
}

// List returns a list of all nvim configurations managed by nvmgr.
func List() ([]string, error) {
	entries, err := os.ReadDir(ConfigDir())
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	var configs []string
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), ConfigPrefix) {
			configs = append(configs, strings.TrimPrefix(e.Name(), ConfigPrefix))
		}
	}

	return configs, nil
}

// CreateConfigsFile creates the configs.json file.
func CreateConfigsFile() (*os.File, error) {
	path := nvmgrConfigsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return os.Create(path)
}

// ReadConfigsFile reads the configs.json file and returns a list of configs.
func ReadConfigs() ([]Config, error) {
	var configs []Config
	path := nvmgrConfigsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return configs, err
	}
	if err = json.Unmarshal(data, &configs); err != nil {
		return configs, err
	}
	return configs, nil
}

// WriteConfigs writes configs to the configs file
func WriteConfigs(configs []Config) error {
	bytes, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configs: %w", err)
	}

	path := nvmgrConfigsPath()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, bytes, 0o644); err != nil {
		return fmt.Errorf("failed to write configs file: %w", err)
	}

	return nil
}

var ErrConfigNotFound = fmt.Errorf("config not found")

// GetConfig retrieves a single config by name.
func GetConfig(name string) (Config, error) {
	configs, err := ReadConfigs()
	if err != nil {
		return Config{}, err
	}
	for _, c := range configs {
		if c.Name == name {
			return c, nil
		}
	}
	return Config{}, ErrConfigNotFound
}

// AddConfig adds a new config to the config file.
func AddConfig(config Config) error {
	configs, err := ReadConfigs()
	if err != nil {
		return err
	}
	configs = append(configs, config)
	if err = WriteConfigs(configs); err != nil {
		return err
	}
	return nil
}
