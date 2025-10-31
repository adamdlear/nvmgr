package cmd

import (
	"bufio"
	"fmt"
	"os"

	metadata "github.com/adamdlear/nvmgr/internal"
	"github.com/adamdlear/nvmgr/internal/configs"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup nvmgr for first time use",
	Long:  `Initialize nvmgr for first time user. This command backs up your existing config to a name of your choice`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\nEnter a name for your existing Neovim configuration (Example: main, work, backup): ")
		reader := bufio.NewReader(os.Stdin)
		name, _ := reader.ReadString('\n')
		fmt.Printf("You entered: %s\n", name)

		err := os.Rename(os.ExpandEnv("$HOME/.config/nvim"), os.ExpandEnv("$HOME/.config/nvim-"+name))
		if err != nil {
			fmt.Printf("Error renaming existing config: %v\n", err)
			return
		}

		// add the metadata file
		err = metadata.Write(os.ExpandEnv("$HOME/.config/nvim-"+name), name, "Initial configuration")
		if err != nil {
			fmt.Printf("Error writing metadata: %v\n", err)
			return
		}

		// create symlink to the newly named config
		err = configs.Activate(name)
		if err != nil {
			fmt.Printf("Error activating config: %v\n", err)
			return
		}

		fmt.Println("nvmgr initialized successfully!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
