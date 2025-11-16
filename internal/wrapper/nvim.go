package wrapper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adamdlear/nvmgr/internal/state"
)

func ExecuteNvim(args []string) error {
	nvimPath, err := findRealNvim()
	if err != nil {
		return err
	}

	state, err := state.LoadState()
	if err != nil {
		return err
	}

	cmd := exec.Command(nvimPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if state.Current != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("NVIM_APPNAME=%s", state.Current))
	} else {
		cmd.Env = os.Environ()
	}

	return cmd.Run()
}

func findRealNvim() (string, error) {
	// Get the current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return "", err
	}
	currentExe, err = filepath.EvalSymlinks(currentExe)
	if err != nil {
		return "", err
	}

	// Get PATH
	pathEnv := os.Getenv("PATH")
	paths := strings.SplitSeq(pathEnv, ":")

	// Search for nvim in PATH, excluding our wrapper
	for dir := range paths {
		nvimPath := filepath.Join(dir, "nvim")

		// Skip if it doesn't exist
		if _, err := os.Stat(nvimPath); os.IsNotExist(err) {
			continue
		}

		// Resolve symlinks
		realPath, err := filepath.EvalSymlinks(nvimPath)
		if err != nil {
			continue
		}

		// Skip if it's our wrapper
		if realPath == currentExe {
			continue
		}

		// Check if it's actually executable
		if info, err := os.Stat(realPath); err == nil {
			if info.Mode()&0o111 != 0 {
				return nvimPath, nil
			}
		}
	}

	// Try common locations as fallback
	commonPaths := []string{
		"/usr/bin/nvim",
		"/usr/local/bin/nvim",
		"/opt/homebrew/bin/nvim",
		"/opt/local/bin/nvim",
	}

	for _, path := range commonPaths {
		realPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			continue
		}
		if realPath == currentExe {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("could not find real nvim binary")
}
