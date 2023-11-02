package cmd

import (
	"context"
	"io"
	"os"

	"github.com/gomicro/concord/report"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.AddCommand(NewApplyTeamsCmd(os.Stdout))
}

func NewApplyTeamsCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "teams",
		Args:              cobra.ExactArgs(1),
		Short:             "Apply a teams configuration",
		Long:              `Apply teams in a configuration against github`,
		PersistentPreRunE: setupClient,
		RunE:              applyTeamsRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyTeamsRun(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	file := args[0]

	org, err := readManifest(file)
	if err != nil {
		return handleError(cmd, err)
	}

	report.PrintHeader("Org")
	report.Println()

	return teamsRun(ctx, cmd, args, org, false)
}
