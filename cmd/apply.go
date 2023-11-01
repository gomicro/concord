package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewApplyCmd(os.Stdout))
}

func NewApplyCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "apply",
		Args:              cobra.ExactArgs(1),
		Short:             "Apply a github configuration",
		Long:              `Apply a github configuration`,
		PersistentPreRunE: setupClient,
		RunE:              applyRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyRun(cmd *cobra.Command, args []string) error {
	return nil
}
