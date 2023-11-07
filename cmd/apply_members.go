package cmd

import (
	"io"
	"os"

	"github.com/gomicro/concord/manifest"
	"github.com/gomicro/concord/report"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.AddCommand(NewApplyMembersCmd(os.Stdout))
}

func NewApplyMembersCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "members",
		Args:  cobra.ExactArgs(1),
		Short: "Apply a members configuration",
		Long:  `Apply members in a configuration against github`,
		RunE:  applyMembersRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyMembersRun(cmd *cobra.Command, args []string) error {
	file := args[0]
	cmd.SetContext(manifest.WithManifest(cmd.Context(), file))

	report.PrintHeader("Org")
	report.Println()

	return membersRun(cmd, args, false)
}
