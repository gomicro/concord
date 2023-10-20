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
		Use:   "apply",
		Short: "Apply a github configuration",
		Long:  `Apply a github configuration`,
		RunE:  applyRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyRun(cmd *cobra.Command, args []string) error {
	return nil
}
