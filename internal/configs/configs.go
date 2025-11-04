package configs

import (
	"os"
	"path/filepath"
	"strings"
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

func List() []string {
	entries, _ := os.ReadDir(ConfigDir())

	var configs []string
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), ConfigPrefix) {
			continue
		}
		configs = append(configs, e.Name())
	}

	return configs
}
