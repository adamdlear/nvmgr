/*
Copyright © 2025 Adam Lear
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	metadata "github.com/adamdlear/nvmgr/internal"
	"github.com/adamdlear/nvmgr/internal/configs"
	"github.com/adamdlear/nvmgr/internal/symlink"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Backs up your existing config",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Enter a name for your existing Neovim config (ex: main, work, backup...): ")
		reader := bufio.NewReader(os.Stdin)
		name, _ := reader.ReadString('\n')

		existingPath := configs.ConfigPath("nvim")
		newPath := configs.ConfigPath(configs.ConfigPrefix + name)

		err := os.Rename(os.ExpandEnv(existingPath), os.ExpandEnv(newPath))
		if err != nil {
			return fmt.Errorf("error renaming exisiting config :%w", err)
		}

		err = metadata.Write(os.ExpandEnv(newPath), name, "Initial Setup")
		if err != nil {
			return fmt.Errorf("error writing metadata: %w", err)
		}

		err = symlink.Activate(name)
		if err != nil {
			return fmt.Errorf("error activating config %q: %w", name, err)
		}

		fmt.Println("nvmgr setup successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
