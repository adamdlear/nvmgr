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

var distributions = map[string]string{
	"astronvim": "https://github.com/AstroNvim/AstroNvim",
	"kickstart": "https://github.com/nvim-lua/kickstart.nvim.git",
	"lazyvim":   "https://github.com/LazyVim/starter",
	"nvchad":    "https://github.com/NvChad/starter",
}

var installCmd = &cobra.Command{
	Use:   "install <distribution-name | git-url> [name]",
	Args:  cobra.RangeArgs(1, 2),
	Short: "Install a config from a git url or a recognized distribution",
	Example: `
# Install a config from a recognized distribution
nvmgr install lazyvim

# Install a config from a git url and give it a name
nvmgr install https://github.com/AstroNvim/AstroNvim astro`,
	RunE: func(cmd *cobra.Command, args []string) error {
		source := args[0]
		url := source

		if repoURL, ok := distributions[source]; ok {
			url = repoURL
		}

		var name string
		if len(args) == 2 {
			name = args[1]
		} else if _, ok := distributions[source]; ok {
			name = source
		} else {
			parts := strings.Split(url, "/")
			name = parts[len(parts)-1]
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
	var longDesc strings.Builder
	longDesc.WriteString("Install a config from a git url or a recognized distribution.\n\n")
	longDesc.WriteString("Recognized distributions:\n")
	for name := range distributions {
		longDesc.WriteString(fmt.Sprintf("  - %s\n", name))
	}
	installCmd.Long = longDesc.String()

	rootCmd.AddCommand(installCmd)
}
