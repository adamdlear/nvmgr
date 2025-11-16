package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/adamdlear/nvmgr/internal/state"
	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Use:   "launch [name]",
	Short: "Launch Neovim with the specified config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		config, err := state.GetConfig(name)
		if err != nil {
			return fmt.Errorf("could not find config %q: %w", name, err)
		}

		env := os.Environ()
		env = append(env, fmt.Sprintf("NVIM_APPNAME=%s", config.Path))

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
