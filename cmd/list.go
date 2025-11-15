package cmd

import (
	"fmt"
	"strings"
	"time"

	metadata "github.com/adamdlear/nvmgr/internal"
	"github.com/adamdlear/nvmgr/internal/configs"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Neovim configurations with metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		names, err := configs.List()
		if err != nil {
			return err
		}

		fmt.Printf("%-15s %-45s %-25s %s\n", "NAME", "PATH", "CREATED", "DESCRIPTION")
		fmt.Println(strings.Repeat("-", 100))

		for _, n := range names {
			dirPath := configs.ConfigDir() + "/" + n
			meta, err := metadata.Read(dirPath)
			if err != nil {
				// fallback if no metadata file exists
				fmt.Printf("%-15s %-45s %-25s %s\n",
					strings.TrimPrefix(n, configs.ConfigPrefix),
					dirPath,
					"unknown",
					"",
				)
				continue
			}
			fmt.Printf("%-15s %-45s %-25s %s\n",
				meta.Name,
				dirPath,
				meta.CreatedAt.Format(time.RFC822),
				meta.Description,
			)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
