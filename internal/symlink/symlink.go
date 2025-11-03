package symlink

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adamdlear/nvmgr/internal/configs"
)

func ActiveLink() string {
	return filepath.Join(configs.ConfigDir(), "nvim")
}

func Activate(name string) error {
	target := configs.ConfigPath(name)
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
