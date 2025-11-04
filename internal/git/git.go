package git

import (
	"fmt"
	"os"
	"os/exec"
)

func Clone(url string, target string) error {
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("could not find git command")
	}

	cmd := exec.Command("git", "clone", "--depth=1", url, target)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
