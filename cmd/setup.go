package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adamdlear/nvmgr/internal/state"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup nvmgr for your machine",
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

		fmt.Println("Successfully setup nvmgr")
		fmt.Printf("View your saved configs with %q\n", "nvmgr list")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
