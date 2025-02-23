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
	applyCmd.AddCommand(NewApplyTeamsCmd(os.Stdout))
}

func NewApplyTeamsCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "teams",
		Short: "Apply a teams configuration",
		Long:  `Apply teams in a configuration against github`,
		RunE:  applyTeamsRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyTeamsRun(cmd *cobra.Command, args []string) error {
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

	err = teamsRun(cmd, args)
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

func teamsRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	org, err := manifest.OrgFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	scrb.BeginDescribe("Teams")
	defer scrb.EndDescribe()

	tms, err := clt.GetTeams(ctx, org.Name)
	if err != nil {
		return handleError(cmd, err)
	}

	missing, managed, unmanaged := getTeamsBreakdown(org.Teams, tms)

	for _, mt := range missing {
		scrb.BeginDescribe(mt)
		scrb.EndDescribe()

		clt.CreateTeam(ctx, org.Name, mt)

		missing, _, _ := getTeamMembersBreakdown(mt, org.People, nil)

		for _, m := range missing {
			clt.InviteTeamMember(ctx, org.GetName(), mt, m)
		}
	}

	for _, mt := range managed {
		scrb.BeginDescribe(mt)
		scrb.EndDescribe()

		scrb.Done("team exists in github")

		ms, err := clt.GetTeamMembers(ctx, org.Name, mt)
		if err != nil {
			return handleError(cmd, err)
		}

		missing, managed, unmanaged := getTeamMembersBreakdown(mt, org.People, ms)
		for _, m := range missing {
			clt.InviteTeamMember(ctx, org.GetName(), mt, m)
		}

		for _, m := range managed {
			scrb.Done(m + " exists in team")
		}

		for _, m := range unmanaged {
			scrb.Done(m + " exists in team but not in manifest")
		}
	}

	for _, mt := range unmanaged {
		scrb.BeginDescribe(mt)
		scrb.EndDescribe()
		scrb.Done("team exists in github but not in manifest") // TODO: warn
	}

	return nil
}

func getTeamsBreakdown(manifest []string, teams []*github.Team) (missing []string, managed []string, unmanaged []string) {
	for _, t := range teams {
		if managedTeam(manifest, t.GetName()) {
			managed = append(managed, t.GetName())
		} else {
			unmanaged = append(unmanaged, t.GetName())
		}
	}

	for _, m := range manifest {
		found := false
		for _, t := range teams {
			if strings.EqualFold(m, t.GetName()) {
				found = true
			}
		}

		if !found {
			missing = append(missing, m)
		}
	}

	return
}

func managedTeam(manifestTeams []string, name string) bool {
	for _, t := range manifestTeams {
		if strings.EqualFold(t, name) {
			return true
		}
	}

	return false
}

func getTeamMembersBreakdown(team string, people []*gh_pb.People, members []*github.User) (missing []string, managed []string, unmanaged []string) {
	for _, m := range members {
		if managedTeamMember(people, m.GetLogin()) {
			managed = append(managed, m.GetLogin())
		} else {
			unmanaged = append(unmanaged, m.GetLogin())
		}
	}

	for _, p := range people {
		found := false
		for _, m := range members {
			if strings.EqualFold(p.Username, m.GetLogin()) {
				found = true
			}
		}

		if !found {
			for _, t := range p.Teams {
				if strings.EqualFold(t, team) {
					missing = append(missing, p.Username)
				}
			}
		}
	}

	return
}

func managedTeamMember(manifestPeople []*gh_pb.People, name string) bool {
	for _, p := range manifestPeople {
		if strings.EqualFold(p.Username, name) {
			return true
		}
	}

	return false
}
