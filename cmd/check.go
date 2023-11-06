package cmd

import (
	"io"
	"os"

	"github.com/gomicro/concord/manifest"
	"github.com/gomicro/concord/report"
	"github.com/spf13/cobra"
)

var checkCmd = NewCheckCmd(os.Stdout)

func init() {
	rootCmd.AddCommand(checkCmd)
}

func NewCheckCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Args:  cobra.ExactArgs(1),
		Short: "Check a github configuration",
		Long:  `Check a configuration against what exists in github`,
		RunE:  checkRun,
	}

	cmd.SetOut(out)

	return cmd
}

func checkRun(cmd *cobra.Command, args []string) error {
	file := args[0]
	cmd.SetContext(manifest.WithManifest(cmd.Context(), file))

	report.PrintHeader("Org")
	report.Println()

	err := membersRun(cmd, args, true)
	if err != nil {
		return handleError(cmd, err)
	}

	err = teamsRun(cmd, args, true)
	if err != nil {
		return handleError(cmd, err)
	}

	err = reposRun(cmd, args, true)
	if err != nil {
		return handleError(cmd, err)
	}

	return nil
}
