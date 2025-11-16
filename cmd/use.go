package cmd

import (
	"fmt"
	"strings"

	"github.com/adamdlear/nvmgr/internal/state"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Set the active Neovim configuration",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		s, _ := state.LoadState()
		var matches []string
		for _, c := range s.Configs {
			if strings.HasPrefix(c.Name, toComplete) {
				matches = append(matches, c.Name)
				fmt.Println(c.Name)
			}
		}
		return matches, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		exists, err := state.ConfigExists(name)
		if err != nil {
			return fmt.Errorf("failed to get config '%s': %w", name, err)
		}
		if !exists {
			return fmt.Errorf("config '%s' does not exist", name)
		}

		s, err := state.LoadState()
		if err != nil {
			return fmt.Errorf("failed to read current configs: %w", err)
		}

		s.Current = name

		if err := state.SaveState(s); err != nil {
			return fmt.Errorf("faild to activate config '%s': %w", name, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
