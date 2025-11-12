/*
Copyright © 2025 Adam Lear
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/adamdlear/nvmgr/internal/configs"
	"github.com/adamdlear/nvmgr/internal/files"
	"github.com/adamdlear/nvmgr/internal/symlink"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore your machine from nvmgr configuration",
	Long: `Use this command if you no longer want to use nvmgr to manage your Neovim configurations.
	
This command does the following things:
- moves the desired configuration to '~/.config/nvim'	
- removes any remaining symlinks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		nvimDir := configs.ConfigDir() + "/nvim"

		reader := bufio.NewReader(os.Stdin)
		confirmed, err := confirm("Are you sure you want to restore your Neovim configurations?", reader)
		if err != nil {
			return err
		}

		if !confirmed {
			fmt.Println("Canceled restoration. Exiting...")
			return nil
		}

		if _, err = os.Stat(nvimDir); err != nil {
			if os.IsNotExist(err) {
				confirmed, err = confirm("Would you like to create the ~/.config/nvim directory?", reader)
				if err != nil {
					return err
				}

				if !confirmed {
					fmt.Println("The ~/.config/nvim directory needs to exist to continue. Exiting...")
				}
			}
		}

		fmt.Print("Which config would you like to be your main config? ")
		name, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading user input")
		}
		name = strings.TrimPrefix(name, "\n")
		path := configs.ConfigPath(name)

		if err := replaceNvimPath(path, nvimDir); err != nil {
			return err
		}
		if err := replaceNvimPath("~/.local/shared/nvim-"+name, "~/.local/shared/nvim"); err != nil {
			return err
		}
		if err := replaceNvimPath("~/.local/state/nvim-"+name, "~/.local/state/nvim"); err != nil {
			return err
		}
		if err := replaceNvimPath("~/.cache/nvim-"+name, "~/.cache/nvim"); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}

func confirm(question string, reader *bufio.Reader) (bool, error) {
	for {
		fmt.Printf("%s [Y/n]: ", question)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.TrimSpace(response)

		switch response {
		case "Y", "y":
			return true, nil
		case "N", "n":
			return false, nil
		default:
			fmt.Printf("Invalid option %q. Please enter Y or n.\n", response)
		}
	}
}

func replaceNvimPath(old string, new string) error {
	if err := files.CopyDir(old, new); err != nil {
		return fmt.Errorf("could not copy files from %q to %q, %w", old, new, err)
	}
	if err := symlink.Update(old, new); err != nil {
		return fmt.Errorf("could not update link from %q to %q, %w", old, new, err)
	}
	return nil
}
