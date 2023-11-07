package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/gomicro/concord/manifest"
	"github.com/gomicro/concord/report"
	"github.com/spf13/cobra"
)

var applyCmd = NewApplyCmd(os.Stdout)

func init() {
	rootCmd.AddCommand(applyCmd)
}

func NewApplyCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply an org configuration",
		Long:  `Apply an org configuration against github`,
		RunE:  applyRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyRun(cmd *cobra.Command, args []string) error {
	file := cmd.Flags().Lookup("file").Value.String()
	cmd.SetContext(manifest.WithManifest(cmd.Context(), file))

	dry := strings.EqualFold(cmd.Flags().Lookup("dry").Value.String(), "true")

	report.PrintHeader("Org")
	report.Println()

	err := membersRun(cmd, args, dry)
	if err != nil {
		return handleError(cmd, err)
	}

	err = teamsRun(cmd, args, dry)
	if err != nil {
		return handleError(cmd, err)
	}

	err = reposRun(cmd, args, false)
	if err != nil {
		return handleError(cmd, err)
	}

	return nil
}
