package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		reader := bufio.NewReader(os.Stdin)
		confirmed, err := confirm("Are you sure you want to restore your Neovim configurations?", reader)
		if err != nil {
			return err
		}

		if !confirmed {
			fmt.Println("Canceled restoration. Exiting...")
			return nil
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		installDir := filepath.Join(home, ".local", "bin")
		wrapperPath := filepath.Join(installDir, "nvim")

		if err := os.Remove(wrapperPath); err != nil {
			return fmt.Errorf("failed to remove nvim binary at '%s'", wrapperPath)
		}

		fmt.Println("Successfully cleaned up Neovim")

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
