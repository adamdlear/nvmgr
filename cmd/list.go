package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/adamdlear/nvmgr/internal/state"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Neovim configurations",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.LoadState()
		if err != nil {
			return fmt.Errorf("failed to read current configs: %w", err)
		}

		current := s.Current
		configs := s.Configs

		if len(configs) == 0 {
			fmt.Println("No configurations found.")
			return nil
		}

		sort.Slice(configs, func(i, j int) bool {
			ci := s.Configs[i]
			cj := s.Configs[j]

			// If i is the current, it goes first
			if ci.Name == current && cj.Name != current {
				return true
			}
			// If j is the current, it goes first
			if cj.Name == current && ci.Name != current {
				return false
			}

			// Otherwise sort normally (by name here)
			return ci.Name < cj.Name
		})

		for _, c := range configs {
			fmt.Printf("> %s", c.Name)
			if c.Name == current {
				fmt.Print(" (Active)\n")
			} else {
				fmt.Printf("\n")
			}
			fmt.Printf("  Path: %s\n", c.Path)
			fmt.Printf("  Created: %s\n\n", c.CreatedAt.Format(time.RFC822))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
