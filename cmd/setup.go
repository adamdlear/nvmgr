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

			active := false
			name := strings.TrimPrefix(e.Name(), "nvim-")
			if e.Name() == "nvim" {
				name = "main"
				active = true
			}
			path := filepath.Join(configs.ConfigDir(), e.Name())
			timestamp := time.Now()

			config := configs.Config{
				Name:      name,
				Path:      path,
				CreatedAt: timestamp,
				Active:    active,
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
