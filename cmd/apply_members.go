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
	applyCmd.AddCommand(NewApplyMembersCmd(os.Stdout))
}

func NewApplyMembersCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "members",
		Short: "Apply a members configuration",
		Long:  `Apply members in a configuration against github`,
		RunE:  applyMembersRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyMembersRun(cmd *cobra.Command, args []string) error {
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

	err = membersRun(cmd, args)
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

func membersRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	org, err := manifest.OrgFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	scrb.BeginDescribe("Members")
	scrb.EndDescribe()

	ms, err := clt.GetMembers(ctx, org.Name)
	if err != nil {
		return handleError(cmd, err)
	}

	missing, managed, unmanaged := getMemberBreakdown(org.People, ms)

	for _, m := range missing {
		clt.InviteMember(ctx, org.Name, m)
		scrb.Done("Invited " + m)
	}

	for _, m := range managed {
		scrb.Done(m + " exists in github")
	}

	for _, m := range unmanaged {
		scrb.Done(m + " exists in github but not in manifest")
	}

	return nil
}

func getMemberBreakdown(people []*gh_pb.People, members []*github.User) (missing []string, managed []string, unmanaged []string) {
	for _, m := range members {
		if managedMember(people, m) {
			managed = append(managed, *m.Login)
		} else {
			unmanaged = append(unmanaged, *m.Login)
		}
	}

	for _, p := range people {
		found := false
		for _, m := range members {
			if strings.EqualFold(p.Username, m.GetLogin()) {
				found = true
				break
			}
		}

		if !found {
			missing = append(missing, p.Username)
		}
	}

	return
}

func managedMember(manifestMembers []*gh_pb.People, member *github.User) bool {
	for _, mm := range manifestMembers {
		if strings.EqualFold(mm.Username, *member.Login) {
			return true
		}
	}

	return false
}
