package cmd

import (
	"io"
	"os"

	"github.com/gomicro/concord/manifest"
	"github.com/gomicro/concord/report"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.AddCommand(NewApplyReposCmd(os.Stdout))
}

func NewApplyReposCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repos",
		Args:  cobra.ExactArgs(1),
		Short: "Apply a repos configuration",
		Long:  `Apply repos in a configuration against github`,
		RunE:  applyReposRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyReposRun(cmd *cobra.Command, args []string) error {
	file := args[0]
	cmd.SetContext(manifest.WithManifest(cmd.Context(), file))

	report.PrintHeader("Org")
	report.Println()

	return reposRun(cmd, args, false)
}
