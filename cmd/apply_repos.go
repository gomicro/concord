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
	"github.com/gomicro/concord/manifest"
	"github.com/gomicro/concord/report"
	"github.com/google/go-github/v56/github"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func init() {
	applyCmd.AddCommand(NewApplyReposCmd(os.Stdout))
}

func NewApplyReposCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "Apply a repos configuration",
		Long:  `Apply repos in a configuration against github`,
		RunE:  applyReposRun,
	}

	cmd.SetOut(out)

	return cmd
}

func applyReposRun(cmd *cobra.Command, args []string) error {
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

	report.PrintHeader("Org")
	report.Println()

	err = reposRun(cmd, args)
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

func reposRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	org, err := manifest.OrgFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	report.Println()
	report.PrintHeader("Repos")
	report.Println()

	for _, r := range org.Repositories {
		report.Println()
		report.PrintHeader(r.Name)
		report.Println()

		err := ensureRepo(ctx, org.Name, r)
		if err != nil {
			return handleError(cmd, err)
		}
	}

	return nil
}

func ensureRepo(ctx context.Context, org string, repo *gh_pb.Repository) error {
	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	ghr, err := clt.GetRepo(ctx, org, repo.Name)
	if err != nil && !errors.Is(err, client.ErrRepoNotFound) {
		return err
	}

	fresh := false
	if errors.Is(err, client.ErrRepoNotFound) {
		clt.CreateRepo(ctx, org, buildRepoState(repo))
		fresh = true
	}

	clt.UpdateRepo(ctx, org, repo.Name, buildRepoEdits(repo, ghr, fresh))

	if len(repo.Labels) > 0 {
		var ghl []string

		if ghr != nil {
			ghl = ghr.Topics
			slices.Sort(ghl)
		}

		l := repo.Labels
		slices.Sort(l)

		if !slices.Equal(ghl, l) {
			clt.SetRepoTopics(ctx, org, repo.Name, l)
		} else {
			report.PrintInfo("labels are [" + strings.Join(l, ", ") + "]")
			report.Println()
		}
	}

	for _, pb := range repo.ProtectedBranches {
		err := setBranchProtection(ctx, org, repo, pb)
		if err != nil {
			return err
		}
	}

	err = setTeamPermissions(ctx, org, repo, ghr)
	if err != nil {
		return err
	}

	err = ensureFiles(ctx, org, repo, ghr)
	if err != nil {
		return err
	}

	return nil
}

func buildRepoEdits(repo *gh_pb.Repository, ghr *github.Repository, fresh bool) *github.Repository {
	edits := &github.Repository{}

	if !fresh && repo.Description != nil && !strings.EqualFold(ghr.GetDescription(), *repo.Description) {
		edits.Description = repo.Description
	}

	if !fresh && repo.Archived != nil && ghr.GetArchived() != *repo.Archived {
		edits.Archived = repo.Archived
	}

	if !fresh && repo.Private != nil && ghr.GetPrivate() != *repo.Private {
		edits.Private = repo.Private
	}

	if !fresh && repo.DefaultBranch != nil && !strings.EqualFold(ghr.GetDefaultBranch(), *repo.DefaultBranch) {
		edits.DefaultBranch = repo.DefaultBranch
	}

	if repo.AutoDeleteHeadBranches != nil && ghr.GetDeleteBranchOnMerge() != *repo.AutoDeleteHeadBranches {
		edits.DeleteBranchOnMerge = repo.AutoDeleteHeadBranches
	}

	if repo.AllowAutoMerge != nil && ghr.GetAllowAutoMerge() != *repo.AllowAutoMerge {
		edits.AllowAutoMerge = repo.AllowAutoMerge
	}

	return edits
}

func buildRepoState(repo *gh_pb.Repository) *github.Repository {
	state := &github.Repository{
		Name: &repo.Name,
	}

	if repo.Description != nil {
		state.Description = repo.Description
	}

	if repo.Archived != nil {
		state.Archived = repo.Archived
	}

	if repo.Private != nil {
		state.Private = repo.Private
	}

	if repo.DefaultBranch != nil {
		state.DefaultBranch = repo.DefaultBranch
	}

	return state
}

func setTeamPermissions(ctx context.Context, org string, repo *gh_pb.Repository, ghr *github.Repository) error {
	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	for p, teams := range repo.Permissions {
		for _, t := range teams.Teams {
			err = clt.AddRepoToTeam(ctx, org, strings.ToLower(t), repo.Name, p)
			if err != nil {
				return err
			}
		}
	}

	// should remove teams without permissions
	managed := map[string]struct{}{}
	for _, ts := range repo.Permissions {
		for _, t := range ts.Teams {
			managed[strings.ToLower(t)] = struct{}{}
		}
	}

	gts, err := clt.GetRepoTeams(ctx, org, repo.Name)
	if err != nil {
		return err
	}

	for _, gt := range gts {
		if _, ok := managed[strings.ToLower(gt.GetName())]; ok {
			continue
		}

		clt.RemoveRepoFromTeam(ctx, org, gt.GetSlug(), repo.Name)
	}

	return nil
}

func ensureFiles(ctx context.Context, org string, repo *gh_pb.Repository, ghr *github.Repository) error {
	// clone down repo
	// copy file to expected location in repo
	// if diff, commit and push PR

	return nil
}

func setBranchProtection(ctx context.Context, org string, repo *gh_pb.Repository, branch *gh_pb.Branch) error {
	state := buildBranchProtectionState(branch)

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	ghpb, err := clt.GetBranchProtection(ctx, org, repo.Name, branch.Name)
	if err != nil && !errors.Is(err, client.ErrBranchProtectionNotFound) {
		return err
	}

	clt.ProtectBranch(ctx, org, repo.Name, branch.Name, state)

	if branch.GetProtection() != nil {
		if ghpb.GetRequiredSignatures().GetEnabled() != branch.GetProtection().GetSignedCommits() {
			clt.RequireSignedCommits(ctx, org, repo.Name, branch.Name)
		} else {
			report.PrintInfo("require signed commits is '" + fmt.Sprintf("%t", branch.GetProtection().GetSignedCommits()) + "'")
			report.Println()
		}
	}

	return nil
}

func buildBranchProtectionState(branch *gh_pb.Branch) *github.ProtectionRequest {
	state := &github.ProtectionRequest{}

	if branch.Protection.RequirePr != nil {
		state.RequiredPullRequestReviews = &github.PullRequestReviewsEnforcementRequest{}
	}

	if branch.Protection.ChecksMustPass != nil {
		state.RequiredStatusChecks = &github.RequiredStatusChecks{
			Checks: []*github.RequiredStatusCheck{},
		}

		if len(branch.Protection.RequiredChecks) > 0 {
			for _, c := range branch.Protection.RequiredChecks {
				state.RequiredStatusChecks.Checks = append(state.RequiredStatusChecks.Checks, &github.RequiredStatusCheck{
					Context: c,
				})
			}
		}
	}

	return state
}
