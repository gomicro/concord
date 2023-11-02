package cmd

import (
	"context"
	"io"
	"os"

	"github.com/gomicro/concord/report"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.AddCommand(NewApplyReposCmd(os.Stdout))
}

func NewApplyReposCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "repos",
		Args:              cobra.ExactArgs(1),
		Short:             "Apply a repos configuration",
		Long:              `Apply repos in a configuration against github`,
		PersistentPreRunE: setupClient,
		RunE:              applyReposRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyReposRun(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	file := args[0]

	org, err := readManifest(file)
	if err != nil {
		return handleError(cmd, err)
	}

	report.PrintHeader("Org")
	report.Println()

	return reposRun(ctx, cmd, args, org, false)
}
