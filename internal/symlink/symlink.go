package symlink

import (
	"fmt"
	"os"
	"path/filepath"
)

func ActiveLink() (string, error) {
	return os.Readlink("nvim")
}

func Activate(target string) error {
	link, err := ActiveLink()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(link); err != nil {
		return fmt.Errorf("failed to remove existing config: %w", err)
	}

	if err := os.Symlink(target, link); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

func ActiveName() (string, error) {
	link, err := ActiveLink()
	if err != nil {
		return "", err
	}

	target, err := os.Readlink(link)
	if err != nil {
		return "", fmt.Errorf("no active config (symlink missing or broken)")
	}
	return filepath.Base(target), nil
}

// Update atomically updates the symlink
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
