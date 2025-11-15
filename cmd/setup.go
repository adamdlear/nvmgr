/*
Copyright © 2025 Adam Lear
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adamdlear/nvmgr/internal/configs"
	"github.com/spf13/cobra"
)

// var setupCmd = &cobra.Command{
// 	Use:   "setup",
// 	Short: "Backs up your existing config",
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		fmt.Println("Enter a name for your existing Neovim config (ex: main, work, backup...): ")
// 		reader := bufio.NewReader(os.Stdin)
// 		name, _ := reader.ReadString('\n')
// 		name = strings.TrimPrefix(name, "\n")
//
// 		existingPath := configs.ConfigPath("nvim")
// 		newPath := configs.ConfigPath(configs.ConfigPrefix + name)
//
// 		err := os.Rename(os.ExpandEnv(existingPath), os.ExpandEnv(newPath))
// 		if err != nil {
// 			return fmt.Errorf("error renaming exisiting config :%w", err)
// 		}
//
// 		err = metadata.Write(os.ExpandEnv(newPath), name, "Initial Setup")
// 		if err != nil {
// 			return fmt.Errorf("error writing metadata: %w", err)
// 		}
//
// 		err = symlink.Activate(name)
// 		if err != nil {
// 			return fmt.Errorf("error activating config %q: %w", name, err)
// 		}
//
// 		fmt.Println("nvmgr setup successfully!")
// 		return nil
// 	},
// }

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup nvmgr for your machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := os.ReadDir(configs.ConfigDir())
		if err != nil {
			return err
		}

		c := []configs.Config{}

		for _, e := range entries {
			if !strings.HasPrefix(e.Name(), "nvim") {
				continue
			}

			name := strings.TrimPrefix(e.Name(), "nvim-")
			if e.Name() == "nvim" {
				name = "main"
			}
			path := filepath.Join(configs.ConfigDir(), e.Name())
			timestamp := time.Now()

			config := configs.Config{
				Name:      name,
				Path:      path,
				CreatedAt: timestamp,
			}

			c = append(c, config)

			fmt.Printf("Saved config for %s\n", path)
		}

		err = configs.WriteConfigs(c)
		if err != nil {
			return err
		}

		fmt.Println("Successfully setup nvmgr")
		fmt.Printf("View your saved configs with %q\n", "nvmgr list")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
