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

	s, err := state.LoadState()
	if err != nil {
		return err
	}

	cmd := exec.Command(nvimPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if s.Current != "" {
		config, err := s.GetConfig(s.Current)
		if err != nil {
			// If the current config isn't found, fall back to default nvim
			cmd.Env = os.Environ()
		} else {
			appName := filepath.Base(config.Path)
			cmd.Env = append(os.Environ(), fmt.Sprintf("NVIM_APPNAME=%s", appName))
		}
	} else {
		cmd.Env = os.Environ()
	}

	return cmd.Run()
}

func findRealNvim() (string, error) {
	currentExe, err := os.Executable()
	if err != nil {
		return "", err
	}
	currentExe, err = filepath.EvalSymlinks(currentExe)
	if err != nil {
		return "", err
	}

	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, ":")

	for _, dir := range paths {
		nvimPath := filepath.Join(dir, "nvim")

		if _, err := os.Stat(nvimPath); os.IsNotExist(err) {
			continue
		}

		realPath, err := filepath.EvalSymlinks(nvimPath)
		if err != nil {
			continue
		}

		if realPath == currentExe {
			continue
		}

		if info, err := os.Stat(realPath); err == nil {
			if info.Mode()&0o111 != 0 {
				return realPath, nil
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

	// Check fallback paths
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
