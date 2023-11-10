package cmd

import (
	"io"
	"os"

	"github.com/gomicro/concord/manifest"
	"github.com/spf13/cobra"
)

var validateCmd = NewValidateCmd(os.Stdout)

func init() {
	rootCmd.AddCommand(validateCmd)
}

func NewValidateCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate an org configuration",
		Long:  `Validate an org configuration file to ensure it is processable by concord.`,
		RunE:  validateRun,
	}

	cmd.SetOut(out)

	return cmd
}

func validateRun(cmd *cobra.Command, args []string) error {
	file := cmd.Flags().Lookup("file").Value.String()

	_, err := manifest.ReadManifest(file)
	if err != nil {
		return handleError(cmd, err)
	}

	return nil
}
