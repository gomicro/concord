package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/gomicro/concord/client"
	gh_pb "github.com/gomicro/concord/github/v1"
	"github.com/gomicro/concord/report"
	"github.com/google/go-github/v56/github"
	"github.com/spf13/cobra"
)

func init() {
	checkCmd.AddCommand(NewCheckTeamsCmd(os.Stdout))
}

func NewCheckTeamsCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "teams",
		Args:  cobra.ExactArgs(1),
		Short: "Check teams exists in an organization",
		Long:  `Check teams in a configuration against what exists in github`,
		RunE:  checkTeamsRun,
	}

	cmd.SetOut(out)

	return cmd
}

func checkTeamsRun(cmd *cobra.Command, args []string) error {
	file := args[0]

	org, err := readManifest(file)
	if err != nil {
		return handleError(cmd, err)
	}

	report.PrintHeader("Org")
	report.Println()

	return teamsRun(cmd, args, org, true)
}

func teamsRun(cmd *cobra.Command, args []string, org *gh_pb.Organization, dry bool) error {
	ctx := cmd.Context()

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	ghOrg, err := clt.GetOrg(ctx, org.Name)
	if err != nil {
		if errors.Is(err, client.ErrOrgNotFound) {
			return errors.New("org does not exist")
		}

		return handleError(cmd, err)
	}

	report.Println()
	report.PrintHeader("Teams")
	report.Println()

	tms, err := clt.GetTeams(ctx, org.Name)
	if err != nil {
		return handleError(cmd, err)
	}

	for _, t := range tms {
		if !managedTeam(org.Teams, t.GetName()) {
			report.PrintWarn(t.GetName() + " exists in github but not in manifest")
		} else {
			report.PrintInfo(t.GetName() + " exists in github")
		}

		report.Println()
	}

	mts := checkTeams(org.Teams, tms)

	err = createTeams(ctx, org.Name, mts, dry)
	if err != nil {
		return handleError(cmd, err)
	}

	report.Println()
	report.PrintHeader("Team Memberships")
	report.Println()

	if dry {
		// fill in missing teams as fakes
		for i := range mts {
			tms = append(tms, &github.Team{
				Name: &mts[i],
				ID:   github.Int64(-1),
			})
		}
	} else {
		tms, err = clt.GetTeams(ctx, org.Name)
		if err != nil {
			return handleError(cmd, err)
		}
	}

	expected := getExpectedTeamMembers(org.People)

	for _, t := range tms {
		report.PrintHeader(t.GetName())
		report.Println()

		ms, err := clt.GetTeamMembers(ctx, ghOrg.GetID(), t.GetID())
		if err != nil {
			return handleError(cmd, err)
		}

		for _, m := range ms {
			if !managedTeamMember(org.People, m.GetLogin()) {
				report.PrintWarn(m.GetLogin() + " exists in team but not in manifest")
			} else {
				report.PrintInfo(m.GetLogin() + " exists in team")
			}

			report.Println()
		}

		err = inviteTeamMembers(ctx, ghOrg, t, checkTeamMembers(expected[strings.ToLower(t.GetName())], ms), dry)
		if err != nil {
			return handleError(cmd, err)
		}
	}

	return nil
}

func managedTeam(manifestTeams []string, name string) bool {
	for _, t := range manifestTeams {
		if strings.EqualFold(t, name) {
			return true
		}
	}

	return false
}

func checkTeams(manifestTeams []string, githubTeams []*github.Team) []string {
	missing := []string{}

	for _, mt := range manifestTeams {
		found := false
		for _, gt := range githubTeams {
			if strings.EqualFold(mt, *gt.Name) {
				found = true
				break
			}
		}

		if !found {
			missing = append(missing, mt)
		}
	}

	return missing
}

func createTeams(ctx context.Context, org string, teams []string, dry bool) error {
	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	for _, t := range teams {
		if dry {
			report.PrintAdd("create team " + t)
			report.Println()
			continue
		}

		err := clt.CreateTeam(ctx, org, t)
		if err != nil {
			return err
		}

		report.PrintSuccess("created team " + t)
		report.Println()
	}

	return nil
}

func getExpectedTeamMembers(people []*gh_pb.People) map[string][]string {
	expected := map[string][]string{}

	for _, p := range people {
		for _, t := range p.Teams {
			expected[strings.ToLower(t)] = append(expected[strings.ToLower(t)], p.Username)
		}
	}

	return expected
}

func checkTeamMembers(expected []string, members []*github.User) []string {
	missing := []string{}

	for _, em := range expected {
		found := false
		for _, gm := range members {
			if strings.EqualFold(em, *gm.Login) {
				found = true
				break
			}
		}

		if !found {
			missing = append(missing, em)
		}
	}

	return missing
}

func managedTeamMember(manifestPeople []*gh_pb.People, name string) bool {
	for _, p := range manifestPeople {
		if strings.EqualFold(p.Username, name) {
			return true
		}
	}

	return false
}

func inviteTeamMembers(ctx context.Context, org *github.Organization, team *github.Team, members []string, dry bool) error {
	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	for _, m := range members {
		if dry {
			report.PrintWarn("invite " + m + " to team " + *team.Name)
			report.Println()
			continue
		}

		err := clt.InviteTeamMember(ctx, org.GetID(), team.GetID(), m)
		if err != nil {
			return err
		}
		report.PrintSuccess("invited " + m + " to team " + *team.Name)
		report.Println()
	}

	return nil
}
