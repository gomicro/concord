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
	"github.com/gomicro/scribe/color"
	"github.com/google/go-github/v56/github"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func init() {
	applyCmd.AddCommand(NewApplyReposCmd(os.Stdout))
}

func NewApplyReposCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repos [repo_names]",
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

func reposRun(cmd *cobra.Command, args []string, free bool) error {
	ctx := cmd.Context()

	org, err := manifest.OrgFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	scrb.BeginDescribe("Repositories")
	defer scrb.EndDescribe()

	repos, err := clt.GetRepos(ctx, org.Name)
	if err != nil {
		return handleError(cmd, err)
	}

	unmanaged := getUnmanagedRepos(org.Repositories, repos)

	targetMap := map[string]struct{}{}
	if len(args) > 0 {
		for _, r := range args {
			targetMap[r] = struct{}{}
		}
	} else {
		for _, r := range org.Repositories {
			targetMap[r.Name] = struct{}{}
		}
	}

	for _, r := range org.Repositories {
		if _, found := targetMap[r.Name]; found {
			scrb.BeginDescribe(r.Name)
			scrb.EndDescribe()

			if r.Archived != nil && *r.Archived {
				scrb.Print("repo is archived, skipping")
				continue
			}

			err := ensureRepo(ctx, org.Name, free, r)
			if err != nil {
				scrb.Print(color.RedFg(err.Error()))
			}
		}
	}

	if len(args) == 0 {
		for _, mr := range unmanaged {
			scrb.BeginDescribe(mr)
			scrb.EndDescribe()
			scrb.Print(color.YellowFg("repo exists in github but not in manifest"))
		}
	}

	return nil
}

func getUnmanagedRepos(manifest []*gh_pb.Repository, repos []*github.Repository) []string {
	managed := []string{}
	for _, r := range manifest {
		managed = append(managed, r.Name)
	}

	unmanaged := []string{}
	for _, r := range repos {
		if !slices.Contains(managed, r.GetName()) {
			unmanaged = append(unmanaged, r.GetName())
		}
	}

	return unmanaged
}

func ensureRepo(ctx context.Context, org string, free bool, repo *gh_pb.Repository) error {
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
		clt.CreateRepo(ctx, scrb, org, buildRepoState(repo))
		clt.InitRepo(ctx, scrb, org, repo.Name, *repo.DefaultBranch)
		fresh = true
	}

	clt.UpdateRepo(ctx, scrb, org, repo.Name, buildRepoEdits(repo, ghr, fresh))

	if len(repo.Labels) > 0 {
		var ghl []string

		if ghr != nil {
			ghl = ghr.Topics
			slices.Sort(ghl)
		}

		l := repo.Labels
		slices.Sort(l)

		if !slices.Equal(ghl, l) {
			clt.SetRepoTopics(ctx, scrb, org, repo.Name, l)
		} else {
			scrb.Print("labels are [" + strings.Join(l, ", ") + "]")
		}
	}

	if free && repo.GetPrivate() {
		scrb.Print("skipping branch protection for private repo on free plan")
	} else {
		for _, pb := range repo.ProtectedBranches {
			err := setBranchProtection(ctx, org, repo, pb)
			if err != nil {
				return err
			}
		}
	}

	// if repo is fresh, we can't do anything with teams yet
	if !fresh {
		err = setTeamPermissions(ctx, org, repo, ghr)
		if err != nil {
			return err
		}
	}

	err = ensureFiles(ctx, org, repo, ghr)
	if err != nil {
		return err
	}

	return nil
}

func buildRepoEdits(repo *gh_pb.Repository, ghr *github.Repository, fresh bool) *github.Repository {
	edits := &github.Repository{}

	if !fresh && repo.Archived != nil {
		if ghr.GetArchived() != *repo.Archived {
			edits.Archived = repo.Archived
		}
		// Nothing else can be done with archived repos
		if *repo.Archived {
			scrb.Print("repo " + repo.Name + " is archived, skipping")
			return edits
		}
	}

	if !fresh && repo.Description != nil && !strings.EqualFold(ghr.GetDescription(), *repo.Description) {
		edits.Description = repo.Description
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
	if len(repo.Permissions) == 0 {
		return nil
	}

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	for p, teams := range repo.Permissions {
		for _, t := range teams.Teams {
			err = clt.AddRepoToTeam(ctx, scrb, org, strings.ToLower(t), repo.Name, p)
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
		return fmt.Errorf("remove unamaged teams: %w", err)
	}

	for _, gt := range gts {
		if _, ok := managed[strings.ToLower(gt.GetName())]; ok {
			continue
		}

		clt.RemoveRepoFromTeam(ctx, scrb, org, gt.GetSlug(), repo.Name)
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

	err = clt.ProtectBranch(ctx, scrb, org, repo.Name, branch.Name, state)
	if err != nil {
		return err
	}

	if branch.GetProtection() != nil {
		err = clt.SetRequireSignedCommits(ctx, scrb, org, repo.Name, branch.Name, branch.GetProtection().GetSignedCommits())
		if err != nil {
			return err
		}
	}

	return nil
}

func buildBranchProtectionState(branch *gh_pb.Branch) *github.ProtectionRequest {
	state := &github.ProtectionRequest{}

	if branch.Protection.RequirePr != nil && *branch.Protection.RequirePr {
		state.RequiredPullRequestReviews = &github.PullRequestReviewsEnforcementRequest{}
	}

	if branch.Protection.ChecksMustPass != nil && *branch.Protection.ChecksMustPass {
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
