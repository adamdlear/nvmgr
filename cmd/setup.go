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
		fmt.Println("Setting up nvmgr...")

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
		execPath, err = filepath.EvalSymlinks(execPath)
		if err != nil {
			return err
		}

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

		// Copy nvmgr to installation directory if not already there
		if err := files.CopyFile(execPath, nvmgrPath); err != nil {
			return fmt.Errorf("failed to install nvmgr: %w", err)
		}
		if err := os.Chmod(nvmgrPath, 0o755); err != nil {
			return err
		}

		// Create symlink for wrapper, removing old one if it exists
		_ = os.Remove(wrapperPath) // Ignore error if it doesn't exist
		if err := os.Symlink(nvmgrPath, wrapperPath); err != nil {
			return fmt.Errorf("failed to create nvim wrapper: %w", err)
		}
		fmt.Printf("nvmgr and nvim wrapper installed in %s\n", installDir)

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

	fmt.Println("nvmgr initialized successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Make sure", installDir, "is in your PATH")
	fmt.Println(`  2. To apply changes, open a new shell or run 'rehash' (zsh) or 'hash -r' (bash).`)
	fmt.Println("  3. Create a new config: nvmgr new myconfig")
	fmt.Println("  4. Switch to it: nvmgr use myconfig")
	fmt.Println("  5. Run nvim as usual!")

	return nil
}
