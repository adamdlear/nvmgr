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
