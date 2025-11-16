package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adamdlear/nvmgr/internal/files"
	"github.com/adamdlear/nvmgr/internal/state"
	"github.com/spf13/cobra"
)

var from string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [name]",
	Args:  cobra.ExactArgs(1),
	Short: "Create a new Neovim config (optionally from another)",
	Example: `# Create a blank config
nvmgr new my-config

# Clone an existing config
nvmgr new my-config --from main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		configPrefix := "nvim-"
		configDir, err := state.GetConfigDir()
		if err != nil {
			return err
		}

		newPath := filepath.Join(configDir, configPrefix+name)

		if _, err := os.Stat(newPath); err == nil {
			return fmt.Errorf("config %q already exists", name)
		}

		if from != "" {
			fromPath := filepath.Join(configDir, configPrefix+name)
			if _, err := os.Stat(fromPath); os.IsNotExist(err) {
				return fmt.Errorf("source config %q not found", from)
			}

			if err := files.CopyDir(fromPath, newPath); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(newPath, 0o755); err != nil {
				return err
			}
		}

		config := &state.Config{
			Name:      name,
			Path:      newPath,
			CreatedAt: time.Now(),
		}
		if err := state.SaveConfig(config); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Created new config: %s (%s)\n", name, newPath)
		return nil
	},
}

func init() {
	newCmd.Flags().StringVarP(&from, "from", "f", "", "clone from an existing config")
	rootCmd.AddCommand(newCmd)
}
