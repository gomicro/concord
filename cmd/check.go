package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	gh_pb "github.com/gomicro/concord/github/v1"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func init() {
	rootCmd.AddCommand(NewCheckCmd(os.Stdout))
}

func NewCheckCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "check",
		Args:              cobra.ExactArgs(1),
		Short:             "Check a github configuration",
		Long:              `Check a configuration against what exists in github`,
		PersistentPreRunE: setupClient,
		RunE:              checkRun,
	}

	cmd.SetOut(out)

	return cmd
}

func checkRun(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	file := args[0]

	org, err := readManifest(file)
	if err != nil {
		return err
	}

	orgExists, err := clt.OrgExists(ctx, org.Name)
	if err != nil {
		return err
	}

	if !orgExists {
		return errors.New("org does not exist")
	}

	// check teams exist
	tms, err := clt.GetTeams(ctx, org.Name)
	if err != nil {
		return err
	}

	mts := checkTeams(ctx, org.Teams, tms)

	err = createTeams(ctx, org.Name, mts, true)
	if err != nil {
		return err
	}

	// fill in missing teams as fakes
	for i := range mts {
		tms = append(tms, &github.Team{
			Name: &mts[i],
			ID:   github.Int64(-1),
		})
	}

	// check people exist
	ps, err := clt.GetMembers(ctx, org.Name)
	if err != nil {
		return err
	}

	err = inviteMembers(ctx, org.Name, checkMembers(ctx, org.People, ps), true)
	if err != nil {
		return err
	}

	em := getExpectedTeamMembers(org.People)

	for _, t := range tms {
		// get teams members
		ms, err := clt.GetTeamMembers(ctx, t.GetID())
		if err != nil {
			return err
		}

		err = inviteTeamMembers(ctx, org.Name, t, checkTeamMembers(ctx, em[strings.ToLower(t.GetName())], ms), true)
		if err != nil {
			return err
		}
	}

	// check repos exist

	return nil
}

func checkTeams(ctx context.Context, manifestTeams []string, githubTeams []*github.Team) []string {
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
	for _, t := range teams {
		if dry {
			fmt.Printf("would create team %s\n", t)
			continue
		}

		/*
			_, err := clt.CreateTeam(ctx, org, t)
			if err != nil {
				return err
			}
		*/
	}

	return nil
}

func checkMembers(ctx context.Context, manifestMembers []*gh_pb.People, githubMembers []*github.User) []*gh_pb.People {
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
	for _, m := range members {
		if dry {
			fmt.Printf("would invite %s\n", m.Name)
			continue
		}

		/*
			_, err := clt.InviteMember(ctx, org, m.Name)
			if err != nil {
				return err
			}
		*/
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

func checkTeamMembers(ctx context.Context, expected []string, members []*github.User) []string {
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

func inviteTeamMembers(ctx context.Context, org string, team *github.Team, members []string, dry bool) error {
	for _, m := range members {
		if dry {
			fmt.Printf("would invite %s to team %s\n", m, team.GetName())
			continue
		}

		/*
			_, err := clt.InviteTeamMember(ctx, org, team, m)
			if err != nil {
				return err
			}
		*/
	}

	return nil
}
