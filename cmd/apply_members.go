package cmd

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/gomicro/concord/client"
	gh_pb "github.com/gomicro/concord/github/v1"
	"github.com/gomicro/concord/manifest"
	"github.com/gomicro/concord/report"
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
	ctx := cmd.Context()

	file := cmd.Flags().Lookup("file").Value.String()
	cmd.SetContext(manifest.WithManifest(cmd.Context(), file))

	dry := strings.EqualFold(cmd.Flags().Lookup("dry").Value.String(), "true")

	report.PrintHeader("Org")
	report.Println()

	org, err := manifest.OrgFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	report.Println()
	report.PrintHeader("Members")
	report.Println()

	// check people exist
	ps, err := clt.GetMembers(ctx, org.Name)
	if err != nil {
		return handleError(cmd, err)
	}

	for _, p := range ps {
		if !managedMember(org.People, p) {
			report.PrintWarn(p.GetLogin() + " exists in github but not in manifest")
		} else {
			report.PrintInfo(p.GetLogin() + " exists in github")
		}

		report.Println()
	}

	err = inviteMembers(ctx, org.Name, missingMembers(org.People, ps), dry)
	if err != nil {
		return handleError(cmd, err)
	}

	return nil
}

func managedMember(manifestMembers []*gh_pb.People, member *github.User) bool {
	for _, mm := range manifestMembers {
		if strings.EqualFold(mm.Username, *member.Login) {
			return true
		}
	}

	return false
}

func missingMembers(manifestMembers []*gh_pb.People, githubMembers []*github.User) []*gh_pb.People {
	missing := []*gh_pb.People{}

	for _, mm := range manifestMembers {
		found := false
		for _, gm := range githubMembers {
			if strings.EqualFold(mm.Username, *gm.Login) {
				found = true
				break
			}
		}

		if !found {
			missing = append(missing, mm)
		}
	}

	return missing
}

func inviteMembers(ctx context.Context, org string, members []*gh_pb.People, dry bool) error {
	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	for _, m := range members {
		if dry {
			report.PrintAdd("invite " + m.Name)
			report.Println()
			continue
		}

		err := clt.InviteMember(ctx, org, m.Name)
		if err != nil {
			return err
		}

		report.PrintAdd("invited " + m.Name)
		report.Println()
	}

	return nil
}
