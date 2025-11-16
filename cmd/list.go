package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/adamdlear/nvmgr/internal/configs"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Neovim configurations",
	RunE: func(cmd *cobra.Command, args []string) error {
		configurations, err := configs.ReadConfigs()
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Run %q to set up nvmgr", "nvmgr list")
				return nil
			}
			return fmt.Errorf("failed to read %s", configs.NvmgrConfigsPath())
		}

		if len(configurations) == 0 {
			fmt.Println("No configurations found.")
			return nil
		}

		for _, config := range configurations {
			fmt.Printf("> %s\n", config.Name)
			fmt.Printf("  Path: %s\n", config.Path)
			fmt.Printf("  Created: %s\n", config.CreatedAt.Format(time.RFC822))
			if config.Active {
				fmt.Print("  Active\n\n")
			} else {
				fmt.Print("\n")
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
