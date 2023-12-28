package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	versionCmd = NewVersionCmd(os.Stdout)

	// Version is the current version of concord, made available for use through
	// out the application.
	Version string
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func NewVersionCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display concord's version",
		Long:  `Display the version of the concord CLI.`,
		RunE:  versionRun,
	}

	cmd.SetOut(out)

	return cmd
}

func versionRun(cmd *cobra.Command, args []string) error {
	if Version == "" {
		fmt.Printf("Concord version dev-local\n")
	} else {
		fmt.Printf("Concord version %v\n", Version)
	}

	return nil
}
