package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/adamdlear/nvmgr/internal/state"
	"github.com/spf13/cobra"
)

var delete bool

var removeCmd = &cobra.Command{
	Use:     "remove",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"rm"},
	Short:   "Remove an existing Neovim config from nvmgr",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		s, err := state.LoadState()
		if err != nil {
			fmt.Printf("failed to load current config state")
			return
		}

		config, err := s.GetConfig(name)
		if err != nil {
			fmt.Printf("failed to find config '%s'", name)
			return
		}

		isCurrent := s.Current == name
		if isCurrent && delete {
			fmt.Printf("cannot delete active config '%s'. Please switch to another config first with 'nvmgr use'", name)
			return
		}

		if delete {
			fmt.Printf("Deleting config files from %s...\n", config.Path)
			if err = os.RemoveAll(config.Path); err != nil {
				fmt.Printf("failed to delete config files for '%s'", name)
				return
			}
		}

		s.Configs = slices.DeleteFunc(s.Configs, func(c state.Config) bool {
			return c.Name == name
		})

		if isCurrent {
			s.Current = ""
		}

		if err = state.SaveState(s); err != nil {
			fmt.Printf("failed to save state after removing config '%s'", name)
			return
		}

		if delete {
			fmt.Printf("Successfully deleted config '%s'.\n", name)
		} else {
			if isCurrent {
				fmt.Printf("Successfully deactivated and removed config '%s'.\n", name)
			} else {
				fmt.Printf("Successfully removed config '%s'.\n", name)
			}
		}
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&delete, "delete", "d", false, "Delete the configuration files")
	rootCmd.AddCommand(removeCmd)
}
