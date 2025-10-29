package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adamdlear/nvmgr/internal/configs"
	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Use:   "launch [name]",
	Short: "Launch Neovim with the specified config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		appName := configs.ConfigPrefix + name
		configDir := filepath.Join(os.Getenv("HOME"), ".config", appName)

		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			return fmt.Errorf("config %q does not exist", name)
		}

		env := os.Environ()
		env = append(env, fmt.Sprintf("NVIM_APPNAME=%s", appName))

		nvim := exec.Command("nvim")
		nvim.Env = env
		nvim.Stdin = os.Stdin
		nvim.Stdout = os.Stdout
		nvim.Stderr = os.Stderr

		return nvim.Run()
	},
}

func init() {
	rootCmd.AddCommand(launchCmd)
}
