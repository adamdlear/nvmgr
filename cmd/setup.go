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

		s := &state.State{
			Current: current,
			Configs: configs,
		}

		err = state.SaveState(s)
		if err != nil {
			return err
		}

		execPath, err := os.Executable()
		if err != nil {
			return err
		}
		fmt.Printf("execPath: %s\n", execPath)
		execPath, err = filepath.EvalSymlinks(execPath)
		if err != nil {
			return err
		}
		fmt.Printf("execPath: %s\n", execPath)

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		installDir := filepath.Join(home, ".local", "bin")

		if err := os.MkdirAll(installDir, 0o755); err != nil {
			return err
		}

		nvmgrPath := filepath.Join(installDir, "nvmgr")
		wrapperPath := filepath.Join(installDir, "nvim")

		fmt.Println("Setting up nvmgr...")

		// Copy nvmgr to installation directory if not already there
		if err := files.CopyFile(execPath, nvmgrPath); err != nil {
			return fmt.Errorf("failed to install nvmgr: %w", err)
		}
		if err := os.Chmod(nvmgrPath, 0o755); err != nil {
			return err
		}
		fmt.Printf("✓ Installed nvmgr to %s\n", nvmgrPath)

		// Create symlink for wrapper, removing old one if it exists
		_ = os.Remove(wrapperPath) // Ignore error if it doesn't exist
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
