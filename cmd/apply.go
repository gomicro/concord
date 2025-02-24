package cmd

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/gomicro/concord/client"
	"github.com/gomicro/concord/manifest"
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

	ctx := cmd.Context()

	org, err := manifest.OrgFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	ghOrg, err := clt.GetOrg(ctx, org.Name)
	if err != nil {
		if !errors.Is(err, client.ErrOrgNotFound) {
			return handleError(cmd, err)
		}

		return handleError(cmd, errors.New("organization does not exist"))
	}

	free := ghOrg.GetPlan().GetName() == "free"

	scrb.BeginDescribe("Organization")
	defer scrb.EndDescribe()

	err = orgRun(cmd, args)
	if err != nil {
		return handleError(cmd, err)
	}

	err = membersRun(cmd, args)
	if err != nil {
		return handleError(cmd, err)
	}

	err = teamsRun(cmd, args)
	if err != nil {
		return handleError(cmd, err)
	}

	err = reposRun(cmd, args, free)
	if err != nil {
		return handleError(cmd, err)
	}

	if !dry {
		if !confirm(cmd, "Apply changes? (y/n): ") {
			return nil
		}

		err = clt.Apply()
		if err != nil {
			return handleError(cmd, err)
		}
	}

	return nil
}
