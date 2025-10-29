package configs

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	ConfigPrefix = "nvim-"
)

func ConfigDir() string {
	return filepath.Join(os.Getenv("HOME"), ".config")
}

func ActiveLink() string {
	return filepath.Join(ConfigDir(), "nvim")
}

func ConfigPath(name string) string {
	return filepath.Join(ConfigDir(), ConfigPrefix+name)
}

func Exists(name string) bool {
	_, err := os.Stat(ConfigPath(name))
	return err == nil
}

func Activate(name string) error {
	target := ConfigPath(name)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf("config %q does not exist", name)
	}

	link := ActiveLink()
	_ = os.Remove(link)

	if err := os.Symlink(target, link); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

func ActiveName() (string, error) {
	link := ActiveLink()
	target, err := os.Readlink(link)
	if err != nil {
		return "", fmt.Errorf("no active config (symlink missing or broken)")
	}
	return filepath.Base(target), nil
}
