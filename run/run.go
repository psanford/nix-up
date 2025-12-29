package run

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	dirFlag   string
	forceFlag bool
)

func Command() *cobra.Command {
	cmd := cobra.Command{
		Use:   "run",
		Short: "Pull and run",
		Run:   runAction,
	}

	cmd.Flags().StringVarP(&dirFlag, "dir", "d", "/etc/nixos/nix-cfg", "Directory of git repository")
	cmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force run even if there are no changes")

	return &cmd
}

func runAction(cmd *cobra.Command, args []string) {

	err := os.Chdir(dirFlag)
	if err != nil {
		log.Fatalf("chdir to %s: %s", dirFlag, err)
	}

	out, err := shell("git", "pull")
	fmt.Println(out)
	if err != nil {
		log.Fatalf("git pull err: %s", err)
	}

	if strings.Contains(out, "Already up to date") && !forceFlag {
		log.Println("No changes")
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("get hostname err: %s", err)
	}

	if hostname == "" {
		log.Fatal("no hostname set")
	}

	_, err = os.Stat(hostname)
	if errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Fatalf("no config dir found for %s at %s/%s", hostname, dirFlag, hostname)
		}
	} else if err != nil {
		log.Fatalf("stat dir %s err: %s", hostname, err)
	}

	err = shellStream("nix-channel", "--update")
	if err != nil {
		log.Fatalf("rebuild err: %s", err)
	}

	err = shellStream("nixos-rebuild", "switch", "-I", fmt.Sprintf("nixos-config=/etc/nixos/nix-cfg/%s/configuration.nix", hostname))
	if err != nil {
		log.Fatalf("rebuild err: %s", err)
	}
}

func shell(name string, arg ...string) (string, error) {
	out, err := exec.Command(name, arg...).CombinedOutput()
	return string(out), err
}

func shellStream(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
