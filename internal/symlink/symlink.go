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

func Activate(target string) error {
	link := ActiveLink()
	if err := os.RemoveAll(link); err != nil {
		return fmt.Errorf("failed to remove existing config: %w", err)
	}

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

// Atomically update the symlink
func Update(old string, new string) error {
	newTmp := new + ".tmp"

	if err := os.Remove(newTmp); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.Symlink(old, newTmp); err != nil {
		return err
	}

	if err := os.Rename(newTmp, new); err != nil {
		return err
	}

	return nil
}
