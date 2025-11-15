package cmd

import (
	"fmt"
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
		names, _ := configs.List()
		var matches []string
		for _, n := range names {
			noPref := strings.TrimPrefix(n, configs.ConfigPrefix)
			if strings.HasPrefix(noPref, toComplete) {
				matches = append(matches, n)
				fmt.Println(n)
			}
		}
		return matches, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if !configs.Exists(name) {
			return fmt.Errorf("config %q does not exist", name)
		}

		if err := symlink.Activate(name); err != nil {
			return err
		}

		fmt.Printf("Now using Neovim config: %s\n", strings.TrimPrefix(configs.ConfigPath(name), configs.ConfigDir()+"/"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
