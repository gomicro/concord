package cmd

import (
	"io"
	"os"

	"github.com/gomicro/concord/report"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var checkCmd = NewCheckCmd(os.Stdout)

func init() {
	rootCmd.AddCommand(checkCmd)
}

func NewCheckCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "check",
		Args:              cobra.ExactArgs(1),
		Short:             "Check a github configuration",
		Long:              `Check a configuration against what exists in github`,
		PersistentPreRunE: setupClient,
		RunE:              checkRun,
	}

	cmd.SetOut(out)

	return cmd
}

func checkRun(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	file := args[0]

	org, err := readManifest(file)
	if err != nil {
		return handleError(cmd, err)
	}

	report.PrintHeader("Org")
	report.Println()

	err = _checkMembersRun(ctx, cmd, args, org)
	if err != nil {
		return handleError(cmd, err)
	}

	err = _checkTeamsRun(ctx, cmd, args, org)
	if err != nil {
		return handleError(cmd, err)
	}

	err = _checkReposRun(ctx, cmd, args, org)
	if err != nil {
		return handleError(cmd, err)
	}

	return nil
}
