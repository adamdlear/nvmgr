package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adamdlear/nvmgr/internal/files"
	"github.com/adamdlear/nvmgr/internal/state"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup nvmgr and install the nvim wrapper",
	Long:  `Set up nvmgr by installing the nvim wrapper binary that intercepts nvim calls.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := state.GetConfigDir()
		if err != nil {
			return fmt.Errorf("failed to read from user's config directory: %w", err)
		}

		entries, err := os.ReadDir(configDir)
		if err != nil {
			return err
		}

		current := ""
		configs := []state.Config{}

		for _, e := range entries {
			if !strings.HasPrefix(e.Name(), "nvim") {
				continue
			}

			name := strings.TrimPrefix(e.Name(), "nvim-")
			if e.Name() == "nvim" {
				name = "main"
				current = "main"
			}
			path := filepath.Join(configDir, e.Name())
			timestamp := time.Now()

			config := state.Config{
				Name:      name,
				Path:      path,
				CreatedAt: timestamp,
			}

			configs = append(configs, config)

			fmt.Printf("Saved config for %s\n", path)
		}

		s := state.State{
			Current: current,
			Configs: configs,
		}

		err = state.SaveState(&s)
		if err != nil {
			return err
		}

		execPath, err := os.Executable()
		if err != nil {
			return err
		}
		execPath, err = filepath.EvalSymlinks(execPath)
		if err != nil {
			return err
		}

		installDir := "/usr/local/bin"
		if !isWritable(installDir) {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			installDir = filepath.Join(home, ".local", "bin")
			if err := os.MkdirAll(installDir, 0o755); err != nil {
				return err
			}
		}

		nvmgrPath := filepath.Join(installDir, "nvmgr")
		wrapperPath := filepath.Join(installDir, "nvim")

		fmt.Println("Setting up nvmgr...")

		// Copy nvmgr to installation directory if not already there
		if execPath != nvmgrPath {
			if err := files.CopyFile(execPath, nvmgrPath); err != nil {
				return fmt.Errorf("failed to install nvmgr: %w", err)
			}
			if err := os.Chmod(nvmgrPath, 0o755); err != nil {
				return err
			}
			fmt.Printf("✓ Installed nvmgr to %s\n", nvmgrPath)
		}

		// Check if wrapper already exists and points to us
		if info, err := os.Lstat(wrapperPath); err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				target, err := os.Readlink(wrapperPath)
				if err == nil && target == nvmgrPath {
					fmt.Printf("✓ Wrapper already installed at %s\n", wrapperPath)
					return finishInit(installDir)
				}
			}
			// Wrapper exists but is not our symlink
			return fmt.Errorf("nvim already exists at %s and is not managed by nvmgr. Please remove it first or install nvmgr to a different location", wrapperPath)
		}

		// Create symlink for wrapper
		if err := os.Symlink(nvmgrPath, wrapperPath); err != nil {
			return fmt.Errorf("failed to create nvim wrapper: %w", err)
		}
		fmt.Printf("✓ Installed nvim wrapper at %s\n", wrapperPath)

		fmt.Println("Successfully setup nvmgr")
		fmt.Printf("View your saved configs with %q\n", "nvmgr list")

		return finishInit(installDir)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func isWritable(path string) bool {
	testFile := filepath.Join(path, ".nvmgr_write_test")
	err := os.WriteFile(testFile, []byte("test"), 0o644)
	if err == nil {
		os.Remove(testFile)
		return true
	}
	return false
}

func finishInit(installDir string) error {
	_, err := state.LoadState()
	if err != nil {
		return err
	}

	fmt.Println("\n✓ nvmgr initialized successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Make sure", installDir, "is in your PATH")
	fmt.Println("  2. Create a new config: nvmgr new myconfig")
	fmt.Println("  3. Switch to it: nvmgr use myconfig")
	fmt.Println("  4. Run nvim as usual!")
	fmt.Println("\nNote: Each config will have its own:")
	fmt.Println("  - Config directory: ~/.config/<name>")
	fmt.Println("  - Data directory: ~/.local/share/<name> (plugins, lazy.nvim, etc.)")
	fmt.Println("  - State directory: ~/.local/state/<name>")
	fmt.Println("  - Cache directory: ~/.cache/<name>")

	return nil
}
