package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	metadata "github.com/adamdlear/nvmgr/internal"
	"github.com/adamdlear/nvmgr/internal/configs"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Interactively switch between Neovim configurations",
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir := configs.ConfigDir()
		entries, err := os.ReadDir(configDir)
		if err != nil {
			return fmt.Errorf("failed to read config dir: %w", err)
		}

		var configNames []string  // folder base names (like nvim-work)
		var displayNames []string // human-readable names

		for _, e := range entries {
			if e.IsDir() && strings.HasPrefix(e.Name(), configs.ConfigPrefix) {
				dirPath := filepath.Join(configDir, e.Name())
				meta, err := metadata.Read(dirPath)
				displayName := strings.TrimPrefix(e.Name(), configs.ConfigPrefix)

				if err == nil && meta.Name != "" {
					displayName = meta.Name
				}

				configNames = append(configNames, strings.TrimPrefix(e.Name(), configs.ConfigPrefix))
				displayNames = append(displayNames, displayName)
			}
		}

		if len(configNames) == 0 {
			return fmt.Errorf("no Neovim configs found")
		}

		active, _ := configs.ActiveName()

		fmt.Println("Available configurations:")
		for i, name := range configNames {
			symbol := " "
			if filepath.Base(active) == configs.ConfigPrefix+name {
				symbol = "*"
			}
			fmt.Printf("[%d] %s %s\n", i+1, symbol, displayNames[i])
		}

		fmt.Print("\nEnter number to switch: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var choice int
		_, err = fmt.Sscanf(input, "%d", &choice)
		if err != nil || choice < 1 || choice > len(configNames) {
			return fmt.Errorf("invalid choice")
		}

		name := configNames[choice-1]
		if err := configs.Activate(name); err != nil {
			return err
		}

		fmt.Printf("Switched to config: %s\n", displayNames[choice-1])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
