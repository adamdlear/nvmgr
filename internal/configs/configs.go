package configs

import (
	"os"
	"path/filepath"
)

const (
	ConfigPrefix = "nvim-"
)

func ConfigDir() string {
	return filepath.Join(os.Getenv("HOME"), ".config")
}

func ConfigPath(name string) string {
	return filepath.Join(ConfigDir(), ConfigPrefix+name)
}

func Exists(name string) bool {
	_, err := os.Stat(ConfigPath(name))
	return err == nil
}
