package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gomicro/concord/client"
	gh_pb "github.com/gomicro/concord/github/v1"
	"github.com/gomicro/concord/report"
	"github.com/google/go-github/v56/github"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func init() {
	checkCmd.AddCommand(NewCheckReposCmd(os.Stdout))
}

func NewCheckReposCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "repos",
		Args:              cobra.ExactArgs(1),
		Short:             "Check repos exists in an organization",
		Long:              `Check repos in a configuration against what exists in github`,
		PersistentPreRunE: setupClient,
		RunE:              checkReposRun,
	}

	cmd.SetOut(out)

	return cmd
}

func checkReposRun(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	file := args[0]

	org, err := readManifest(file)
	if err != nil {
		return handleError(cmd, err)
	}

	report.PrintHeader("Org")
	report.Println()

	return reposRun(ctx, cmd, args, org, true)
}

func reposRun(ctx context.Context, cmd *cobra.Command, args []string, org *gh_pb.Organization, dry bool) error {
	report.Println()
	report.PrintHeader("Repos")
	report.Println()

	// ensure all the repos
	for _, r := range org.Repositories {
		report.Println()
		report.PrintHeader(r.Name)
		report.Println()
		err := ensureRepo(ctx, org.Name, r, dry)
		if err != nil {
			return handleError(cmd, err)
		}
	}

	return nil
}

func ensureRepo(ctx context.Context, org string, repo *gh_pb.Repository, dry bool) error {
	ghr, err := clt.GetRepo(ctx, org, repo.Name)
	if err != nil && !errors.Is(err, client.ErrRepoNotFound) {
		return err
	}

	if errors.Is(err, client.ErrRepoNotFound) {
		err = createRepo(ctx, org, repo, dry)
		if err != nil {
			return err
		}
	}

	edits := &github.Repository{}

	if repo.Description != nil && !strings.EqualFold(ghr.GetDescription(), *repo.Description) {
		edits.Description = repo.Description
	}

	if repo.Archived != nil && ghr.GetArchived() != *repo.Archived {
		edits.Archived = repo.Archived
	}

	if repo.Private != nil && ghr.GetPrivate() != *repo.Private {
		edits.Private = repo.Private
	}

	if repo.DefaultBranch != nil && !strings.EqualFold(ghr.GetDefaultBranch(), *repo.DefaultBranch) {
		edits.DefaultBranch = repo.DefaultBranch
	}

	if dry {
		if edits.Description != nil {
			report.PrintAdd("updating description to '" + *edits.Description + "'")
			report.Println()
		}

		if edits.Archived != nil {
			report.PrintAdd("updating archived to '" + fmt.Sprintf("%t", *edits.Archived) + "'")
			report.Println()
		}

		if edits.Private != nil {
			report.PrintAdd("updating private to '" + fmt.Sprintf("%t", *edits.Private) + "'")
			report.Println()
		}

		if edits.DefaultBranch != nil {
			report.PrintAdd("updating default branch to '" + *edits.DefaultBranch + "'")
			report.Println()
		}
	} else {
		err = clt.UpdateRepo(ctx, org, repo.Name, edits)
		if err != nil {
			return err
		}

		if edits.Description != nil {
			report.PrintAdd("updated description to '" + *edits.Description + "'")
			report.Println()
		}

		if edits.Archived != nil {
			report.PrintAdd("updated archived to '" + fmt.Sprintf("%t", *edits.Archived) + "'")
			report.Println()
		}

		if edits.Private != nil {
			report.PrintAdd("updated private to '" + fmt.Sprintf("%t", *edits.Private) + "'")
			report.Println()
		}

		if edits.DefaultBranch != nil {
			report.PrintAdd("updated default branch to '" + *edits.DefaultBranch + "'")
			report.Println()
		}
	}

	err = ensureTopics(ctx, org, repo, ghr, dry)
	if err != nil {
		return err
	}

	/*
		// files
		err = ensureFiles(ctx, org, repo, r, creating, dry)
		if err != nil {
			return err
		}

		// protected branches
		err = ensureProtectedBranches(ctx, org, repo, r, creating, dry)
		if err != nil {
			return err
		}
	*/

	return nil
}

func ensureTopics(ctx context.Context, org string, repo *gh_pb.Repository, ghr *github.Repository, dry bool) error {
	if len(repo.Labels) == 0 {
		return nil
	}

	ghl := ghr.Topics
	slices.Sort(ghl)

	l := repo.Labels
	slices.Sort(l)

	if !slices.Equal(ghl, l) {
		if dry {
			report.PrintAdd("updating topics to [" + strings.Join(l, ", ") + "]")
			report.Println()

			return nil
		}

		err := clt.SetRepoTopics(ctx, org, repo.Name, l)
		if err != nil {
			return err
		}

		report.PrintAdd("updated topics to [" + strings.Join(l, ", ") + "]")
		report.Println()
	} else {
		report.PrintInfo("topics are [" + strings.Join(l, ", ") + "]")
		report.Println()
	}

	return nil
}

func createRepo(ctx context.Context, org string, repo *gh_pb.Repository, dry bool) error {
	state := &github.Repository{
		Name: &repo.Name,
	}

	if repo.Description != nil {
		state.Description = repo.Description
	}

	if repo.Archived != nil {
		state.Archived = repo.Archived
	}

	if len(repo.Labels) > 0 {
		state.Topics = repo.Labels
	}

	if repo.Private != nil {
		state.Private = repo.Private
	}

	if repo.DefaultBranch != nil {
		state.DefaultBranch = repo.DefaultBranch
	}

	if dry {
		report.PrintWarn("creating repo " + repo.Name)
		report.Println()

		if state.Description != nil {
			report.PrintAdd("setting description to '" + *state.Description + "'")
			report.Println()
		}

		if state.Archived != nil {
			report.PrintAdd("setting archived to '" + fmt.Sprintf("%t", *state.Archived) + "'")
			report.Println()
		}

		if len(state.Topics) > 0 {
			report.PrintAdd("setting topics to [" + strings.Join(state.Topics, ", ") + "]")
			report.Println()
		}

		if state.Private != nil {
			report.PrintAdd("setting private to '" + fmt.Sprintf("%t", *state.Private) + "'")
			report.Println()
		}

		if state.DefaultBranch != nil {
			report.PrintAdd("setting default branch to '" + *state.DefaultBranch + "'")
			report.Println()
		}
	} else {
		err := clt.CreateRepo(ctx, org, state)
		if err != nil {
			return err
		}

		report.PrintWarn("created repo " + repo.Name)
		report.Println()

		if state.Description != nil {
			report.PrintAdd("set description to '" + *state.Description + "'")
			report.Println()
		}

		if state.Archived != nil {
			report.PrintAdd("set archived to '" + fmt.Sprintf("%t", *state.Archived) + "'")
			report.Println()
		}

		if len(state.Topics) > 0 {
			report.PrintAdd("set topics to [" + strings.Join(state.Topics, ", ") + "]")
			report.Println()
		}

		if state.Private != nil {
			report.PrintAdd("set private to '" + fmt.Sprintf("%t", *state.Private) + "'")
			report.Println()
		}

		if state.DefaultBranch != nil {
			report.PrintAdd("set default branch to '" + *state.DefaultBranch + "'")
			report.Println()
		}
	}

	return nil
}

func ensureFiles(ctx context.Context, org string, repo *gh_pb.Repository, r *github.Repository, creating, dry bool) error {
	return nil
}

func ensureProtectedBranches(ctx context.Context, org string, repo *gh_pb.Repository, r *github.Repository, creating, dry bool) error {
	if creating && dry {
		for _, pb := range repo.ProtectedBranches {
			report.PrintWarn("create protected branch " + pb.Name + " for repo " + repo.Name)
			report.Println()
			continue
		}

		return nil
	}

	// check wanted protected branches
	for _, pb := range repo.ProtectedBranches {
		_, err := clt.GetBranchProtection(ctx, org, repo.Name, pb.Name)
		if err != nil {
			if errors.Is(err, client.ErrBranchProtectionNotFound) {
				if dry {
					report.PrintWarn("create protected branch " + pb.Name + " for repo " + repo.Name)
					report.Println()
					continue
				}

				/*
					_, err := clt.ProtectBranch(ctx, repo.Name, pb.Name)
					if err != nil {
						return handleError(cmd, err)
						continue
					}
				*/
			}

			return err
		}

		// TODO: Update existing protections
		// ensure require pr
		// ensure checks must pass
		// ensure required checks
		// ensure signed commits
	}

	bs, err := clt.GetBranches(ctx, org, repo.Name)
	if err != nil {
		return err
	}

	// remove unwanted protected branches
	for _, b := range bs {
		found, err := clt.IsBranchProtected(ctx, org, repo.Name, b.GetName())
		if err != nil {
			return err
		}

		if found {
			if dry {
				report.PrintHeader("delete protected branch " + b.GetName() + " for repo " + repo.Name)
				report.Println()
				continue
			}

			/*
				_, err := clt.RemoveBranchProtection(ctx, repo.Name, *p.Name)
				if err != nil {
					return handleError(cmd, err)
				}
			*/
		}
	}

	return nil
}
