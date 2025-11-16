package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adamdlear/nvmgr/cmd"
	"github.com/adamdlear/nvmgr/internal/wrapper"
)

func main() {
	execName := filepath.Base(os.Args[0])

	if execName == "nvim" {
		if err := wrapper.ExecuteNvim(os.Args[1:]); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	cmd.Execute()
}
