package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewCheckCmd(os.Stdout))
}

func NewCheckCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "check",
		Short:             "Check a github configuration",
		Long:              `Check a configuration against what exists in github`,
		PersistentPreRunE: setupClient,
		RunE:              checkRun,
	}

	cmd.SetOut(out)

	return cmd
}

func checkRun(cmd *cobra.Command, args []string) error {
	org, err := readManifest("interxfi.yaml")
	if err != nil {
		return err
	}

	return nil
}
