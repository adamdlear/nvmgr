/*
Copyright © 2025 Adam Lear
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/adamdlear/nvmgr/internal/configs"
	"github.com/adamdlear/nvmgr/internal/git"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Args:  cobra.RangeArgs(1, 2),
	Short: "Install an existing config from a git url",
	Example: `
# Install AstroNvim
nvmgr install https://github.com/AstroNvim/AstroNvim astro`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		var name string

		if len(args) == 1 {
			parts := strings.Split(url, "/")
			name = parts[len(parts)-1]
		} else {
			name = args[1]
		}

		fmt.Printf("Installing %q from %s\n", name, url)

		target := configs.ConfigPath(name)
		err := git.Clone(url, target)
		if err != nil {
			return fmt.Errorf("failed to clone repo %s to dir %s", url, target)
		}

		os.RemoveAll(target + "/.git")

		fmt.Printf("Activate your new config with: nvmgr use %s\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
