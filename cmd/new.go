/*
Copyright © 2025 Adam Lear
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	metadata "github.com/adamdlear/nvmgr/internal"
	"github.com/adamdlear/nvmgr/internal/configs"
	"github.com/spf13/cobra"
)

var (
	from string
	desc string
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [name]",
	Args:  cobra.ExactArgs(1),
	Short: "Create a new Neovim config (optionally from another)",
	Example: `# Create a blank config
nvmgr new my-config

# Clone an existing config
nvmgr new my-config --from main --desc "Experimenting with LSP"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		newPath := filepath.Join(configs.ConfigDir(), configs.ConfigPrefix+name)

		if _, err := os.Stat(newPath); err == nil {
			return fmt.Errorf("config %q already exists", name)
		}

		if from != "" {
			fromPath := filepath.Join(configs.ConfigDir(), configs.ConfigPrefix+name)
			if _, err := os.Stat(fromPath); os.IsNotExist(err) {
				return fmt.Errorf("source config %q not found", from)
			}

			if err := copyDir(fromPath, newPath); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(newPath, 0o755); err != nil {
				return err
			}
		}

		if err := metadata.Write(newPath, name, desc); err != nil {
			return fmt.Errorf("failed to write metadata: %w", err)
		}

		fmt.Printf("Created new config: %s (%s)\n", name, newPath)
		return nil
	},
}

func init() {
	newCmd.Flags().StringVarP(&from, "from", "f", "", "clone from an existing config")
	newCmd.Flags().StringVarP(&desc, "desc", "d", "", "add a short description")
	rootCmd.AddCommand(newCmd)
}

func copyDir(src string, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dest, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(target, data, info.Mode())
	})
}
