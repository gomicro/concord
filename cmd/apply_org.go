package cmd

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/gomicro/concord/client"
	gh_pb "github.com/gomicro/concord/github/v1"
	"github.com/gomicro/concord/manifest"
	"github.com/google/go-github/v56/github"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.AddCommand(NewApplyOrgCmd(os.Stdout))
}

func NewApplyOrgCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "Apply specific org level configuration",
		Long:  `Apply specific org level configuration against github`,
		RunE:  applyOrgRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyOrgRun(cmd *cobra.Command, args []string) error {
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

	exists, err := clt.OrgExists(ctx, org.Name)
	if err != nil {
		return handleError(cmd, err)
	}

	if !exists {
		return handleError(cmd, errors.New("organization does not exist"))
	}

	scrb.BeginDescribe("Organization")
	defer scrb.EndDescribe()

	err = orgRun(cmd, args)
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

func orgRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	org, err := manifest.OrgFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	scrb.BeginDescribe("Permissions")
	scrb.EndDescribe()

	// TODO: this should be broken into two parts, determine if the org exists
	// and then apply the permissions
	err = clt.SetOrgPrivileges(ctx, org.Name, buildOrgState(org))
	if err != nil {
		return handleError(cmd, err)
	}

	return nil
}

func buildOrgState(org *gh_pb.Organization) *github.Organization {
	state := &github.Organization{}

	if org.Permissions != nil {
		if org.Permissions.BasePermissions != nil {
			state.DefaultRepoPermission = org.Permissions.BasePermissions
		}

		if org.Permissions.CreatePrivateRepos != nil {
			state.MembersCanCreatePrivateRepos = org.Permissions.CreatePrivateRepos
		}

		if org.Permissions.CreatePublicRepos != nil {
			state.MembersCanCreatePublicRepos = org.Permissions.CreatePublicRepos
		}
	}

	return state
}
