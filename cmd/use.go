package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adamdlear/nvmgr/internal/configs"
	"github.com/adamdlear/nvmgr/internal/symlink"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Set the active Neovim configuration",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		entries, _ := configs.ReadConfigs()
		var matches []string
		for _, c := range entries {
			if strings.HasPrefix(c.Name, toComplete) {
				matches = append(matches, c.Name)
				fmt.Println(c.Name)
			}
		}
		return matches, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		config, err := configs.GetConfig(name)
		if err != nil {
			return fmt.Errorf("failed to find config %q", name)
		}

		if err = symlink.Activate(config.Path); err != nil {
			return err
		}

		fmt.Printf("Now using Neovim config: %s\n", config.Name)

		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}

		dirsToClean := []string{
			filepath.Join(home, ".local", "share", "nvim"),
			filepath.Join(home, ".local", "state", "nvim"),
			filepath.Join(home, ".cache", "nvim"),
		}

		fmt.Println("Clearing Neovim data, state, and cache directories...")
		for _, dir := range dirsToClean {
			if err := os.RemoveAll(dir); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to remove directory %s: %v\n", dir, err)
			}
		}
		fmt.Println("Cleanup complete. Please reinstall plugins in Neovim.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
