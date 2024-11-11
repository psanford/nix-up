package run

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	dirFlag string
)

func Command() *cobra.Command {
	cmd := cobra.Command{
		Use:   "run",
		Short: "Pull and run",
		Run:   runAction,
	}

	cmd.Flags().StringVarP(&dirFlag, "dir", "d", "/etc/nixos/nix-cfg", "Directory of git repository")

	return &cmd
}

func runAction(cmd *cobra.Command, args []string) {

	err := os.Chdir(dirFlag)
	if err != nil {
		log.Fatalf("chdir to %s: %s", dirFlag, err)
	}

	_, err = shell("git", "pull")
	if err != nil {
		log.Fatalf("git pull err: %s", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("get hostname err: %s", err)
	}

	_, err = os.Stat(hostname)
	if errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Fatalf("no config dir found for %s at %s/%s", hostname, dirFlag, hostname)
		}
	} else if err != nil {
		log.Fatalf("stat dir %s err: %s", hostname, err)
	}

	out, err := shell("nixos-rebuild", "switch", "-I", fmt.Sprintf("nixos-config=/etc/nixos/nix-cfg/%s/configuration.nix", hostname))
	fmt.Println(out)
	if err != nil {
		log.Fatalf("rebuild err: %s", err)
	}
}

func shell(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).CombinedOutput()
}
