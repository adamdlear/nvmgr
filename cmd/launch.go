package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adamdlear/nvmgr/internal/state"
	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Use:   "launch [name]",
	Short: "Launch Neovim with the specified config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		s, err := state.LoadState()
		if err != nil {
			return fmt.Errorf("failed to load state: %w", err)
		}
		config, err := s.GetConfig(name)
		if err != nil {
			return fmt.Errorf("could not find config %q: %w", name, err)
		}

		env := os.Environ()
		appName := filepath.Base(config.Path)
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
