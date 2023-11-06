package cmd

import (
	"io"
	"os"

	"github.com/gomicro/concord/report"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.AddCommand(NewApplyMembersCmd(os.Stdout))
}

func NewApplyMembersCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "members",
		Args:              cobra.ExactArgs(1),
		Short:             "Apply a members configuration",
		Long:              `Apply members in a configuration against github`,
		PersistentPreRunE: setupClient,
		RunE:              applyMembersRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyMembersRun(cmd *cobra.Command, args []string) error {
	file := args[0]

	org, err := readManifest(file)
	if err != nil {
		return handleError(cmd, err)
	}

	report.PrintHeader("Org")
	report.Println()

	return membersRun(cmd, args, org, false)
}
