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

	report.PrintHeader("Org")
	report.Println()

	err := reposRun(cmd, args, dry)
	if err != nil {
		return handleError(cmd, err)
	}

	return nil
}

func reposRun(cmd *cobra.Command, args []string, dry bool) error {
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

		err := ensureRepo(ctx, org.Name, r, dry)
		if err != nil {
			return handleError(cmd, err)
		}
	}

	return nil
}

func ensureRepo(ctx context.Context, org string, repo *gh_pb.Repository, dry bool) error {
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
		err = createRepo(ctx, org, repo, dry)
		if err != nil {
			return err
		}

		fresh = true
	}

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

		if edits.DeleteBranchOnMerge != nil {
			report.PrintAdd("updating auto delete head branches to '" + fmt.Sprintf("%t", *edits.DeleteBranchOnMerge) + "'")
			report.Println()
		}

		if edits.AllowAutoMerge != nil {
			report.PrintAdd("updating allow auto merge to '" + fmt.Sprintf("%t", *edits.AllowAutoMerge) + "'")
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

		if edits.DeleteBranchOnMerge != nil {
			report.PrintAdd("updated auto delete head branches to '" + fmt.Sprintf("%t", *edits.DeleteBranchOnMerge) + "'")
			report.Println()
		}

		if edits.AllowAutoMerge != nil {
			report.PrintAdd("updated allow auto merge to '" + fmt.Sprintf("%t", *edits.AllowAutoMerge) + "'")
			report.Println()
		}
	}

	err = ensureTopics(ctx, org, repo, ghr, dry)
	if err != nil {
		return err
	}

	err = ensureProtectedBranches(ctx, org, repo, ghr, dry)
	if err != nil {
		return err
	}

	err = ensurePermissions(ctx, org, repo, ghr, dry)
	if err != nil {
		return err
	}

	err = ensureFiles(ctx, org, repo, ghr, dry)
	if err != nil {
		return err
	}

	return nil
}

func ensurePermissions(ctx context.Context, org string, repo *gh_pb.Repository, ghr *github.Repository, dry bool) error {
	if len(repo.Permissions) == 0 {
		return nil
	}

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	gts, err := clt.GetRepoTeams(ctx, org, repo.Name)
	if err != nil {
		return err
	}

	gtps := map[string]string{}
	for _, gt := range gts {
		p := gt.GetPermission()
		switch p {
		case "pull":
			gtps[strings.ToLower(gt.GetName())] = "read"
		case "push":
			gtps[strings.ToLower(gt.GetName())] = "write"
		default:
			gtps[strings.ToLower(gt.GetName())] = p
		}
	}

	for p, teams := range repo.Permissions {
		for _, t := range teams.Teams {
			tp, ok := gtps[strings.ToLower(t)]

			if ok && strings.EqualFold(tp, p) {
				report.PrintInfo("team '" + t + "' has permission '" + p + "'")
				report.Println()
				continue
			}

			if dry {
				if ok {
					report.PrintAdd("updating repo to team '" + t + "' with '" + p + "'")
				} else {
					report.PrintAdd("adding repo to team '" + t + "' with '" + p + "'")
					report.Println()
				}
			} else {
				err := clt.AddRepoToTeam(ctx, org, strings.ToLower(t), repo.Name, p)
				if err != nil {
					return err
				}

				if ok {
					report.PrintAdd("updated repo to team '" + t + "' with '" + p + "'")
				} else {
					report.PrintAdd("added repo to team '" + t + "' with '" + p + "'")
					report.Println()
				}
			}
		}
	}

	// should remove teams without permissions
	manTeams := map[string]struct{}{}
	for _, ts := range repo.Permissions {
		for _, t := range ts.Teams {
			manTeams[strings.ToLower(t)] = struct{}{}
		}
	}

	for _, gt := range gts {
		if _, ok := manTeams[strings.ToLower(gt.GetName())]; ok {
			continue
		}

		if dry {
			report.PrintAdd("removing repo from team '" + gt.GetName() + "'")
			report.Println()
		} else {
			err := clt.RemoveRepoFromTeam(ctx, org, gt.GetSlug(), repo.Name)
			if err != nil {
				return err
			}

			report.PrintAdd("removed repo from team '" + gt.GetName() + "'")
			report.Println()
		}
	}

	return nil
}

func ensureTopics(ctx context.Context, org string, repo *gh_pb.Repository, ghr *github.Repository, dry bool) error {
	if len(repo.Labels) == 0 {
		return nil
	}

	var ghl []string

	if ghr != nil {
		ghl = ghr.Topics
		slices.Sort(ghl)
	}

	l := repo.Labels
	slices.Sort(l)

	if !slices.Equal(ghl, l) {
		if dry {
			report.PrintAdd("updating labels to [" + strings.Join(l, ", ") + "]")
			report.Println()

			return nil
		}

		clt, err := client.ClientFromContext(ctx)
		if err != nil {
			return err
		}

		err = clt.SetRepoTopics(ctx, org, repo.Name, l)
		if err != nil {
			return err
		}

		report.PrintAdd("updated labels to [" + strings.Join(l, ", ") + "]")
		report.Println()
	} else {
		report.PrintInfo("labels are [" + strings.Join(l, ", ") + "]")
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

		if state.Private != nil {
			report.PrintAdd("setting private to '" + fmt.Sprintf("%t", *state.Private) + "'")
			report.Println()
		}

		if state.DefaultBranch != nil {
			report.PrintAdd("setting default branch to '" + *state.DefaultBranch + "'")
			report.Println()
		}
	} else {
		clt, err := client.ClientFromContext(ctx)
		if err != nil {
			return err
		}

		err = clt.CreateRepo(ctx, org, state)
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

func ensureFiles(ctx context.Context, org string, repo *gh_pb.Repository, ghr *github.Repository, dry bool) error {
	// clone down repo
	// copy file to expected location in repo
	// if diff, commit and push PR

	return nil
}

func ensureProtectedBranches(ctx context.Context, org string, repo *gh_pb.Repository, ghr *github.Repository, dry bool) error {
	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	for _, pb := range repo.ProtectedBranches {
		_, err := clt.GetBranchProtection(ctx, org, repo.Name, pb.Name)
		if err != nil {
			if errors.Is(err, client.ErrBranchProtectionNotFound) {
				err := createProtectedBranch(ctx, org, repo, pb, dry)
				if err != nil {
					return err
				}

				continue
			}

			return err
		}

		err = UpdateBranchProtection(ctx, org, repo, pb, dry)
		if err != nil {
			return err
		}
	}

	return nil
}

func createProtectedBranch(ctx context.Context, org string, repo *gh_pb.Repository, branch *gh_pb.Branch, dry bool) error {
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

	if dry {
		report.PrintAdd("create protected branch " + branch.Name + " for repo " + repo.Name)
		report.Println()

		if state.RequiredPullRequestReviews != nil {
			report.PrintAdd("setting require pr to '" + fmt.Sprintf("%t", *branch.Protection.RequirePr) + "'")
			report.Println()
		}

		if state.RequiredStatusChecks != nil {
			report.PrintAdd("setting require status checks to '" + fmt.Sprintf("%t", *branch.Protection.ChecksMustPass) + "'")
			report.Println()

			if len(state.RequiredStatusChecks.Checks) > 0 {
				report.PrintAdd("setting required checks to [" + strings.Join(branch.Protection.RequiredChecks, ", ") + "]")
				report.Println()
			}
		}

		err := ensureSignedCommits(ctx, org, repo, branch, dry)
		if err != nil {
			return err
		}

		return nil
	}

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	err = clt.ProtectBranch(ctx, org, repo.Name, branch.Name, state)
	if err != nil {
		return err
	}

	report.PrintWarn("created protected branch " + branch.Name + " for repo " + repo.Name)
	report.Println()

	if state.RequiredPullRequestReviews != nil {
		report.PrintAdd("set require pr to '" + fmt.Sprintf("%t", *branch.Protection.RequirePr) + "'")
		report.Println()
	}

	if state.RequiredStatusChecks != nil {
		report.PrintAdd("set require status checks to '" + fmt.Sprintf("%t", *branch.Protection.ChecksMustPass) + "'")
		report.Println()

		if len(state.RequiredStatusChecks.Checks) > 0 {
			report.PrintAdd("set required checks to [" + strings.Join(branch.Protection.RequiredChecks, ", ") + "]")
			report.Println()
		}
	}

	err = ensureSignedCommits(ctx, org, repo, branch, dry)
	if err != nil {
		return err
	}

	return nil
}

func UpdateBranchProtection(ctx context.Context, org string, repo *gh_pb.Repository, branch *gh_pb.Branch, dry bool) error {
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

	report.PrintInfo("protected branch '" + branch.Name + "' for repo " + repo.Name)
	report.Println()

	if dry {
		if state.RequiredPullRequestReviews != nil {
			report.PrintAdd("updating require pr to '" + fmt.Sprintf("%t", *branch.Protection.RequirePr) + "'")
			report.Println()
		}

		if state.RequiredStatusChecks != nil {
			report.PrintAdd("updating require status checks to '" + fmt.Sprintf("%t", *branch.Protection.ChecksMustPass) + "'")
			report.Println()

			if len(state.RequiredStatusChecks.Checks) > 0 {
				report.PrintAdd("updating required checks to [" + strings.Join(branch.Protection.RequiredChecks, ", ") + "]")
				report.Println()
			}
		}

		err := ensureSignedCommits(ctx, org, repo, branch, dry)
		if err != nil {
			return err
		}

		return nil
	}

	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	err = clt.ProtectBranch(ctx, org, repo.Name, branch.Name, state)
	if err != nil {
		return err
	}

	if state.RequiredPullRequestReviews != nil {
		report.PrintAdd("updated require pr to '" + fmt.Sprintf("%t", *branch.Protection.RequirePr) + "'")
		report.Println()
	}

	if state.RequiredStatusChecks != nil {
		report.PrintAdd("updated require status checks to '" + fmt.Sprintf("%t", *branch.Protection.ChecksMustPass) + "'")
		report.Println()

		if len(state.RequiredStatusChecks.Checks) > 0 {
			report.PrintAdd("updated required checks to [" + strings.Join(branch.Protection.RequiredChecks, ", ") + "]")
			report.Println()
		}
	}

	err = ensureSignedCommits(ctx, org, repo, branch, dry)
	if err != nil {
		return err
	}

	return nil
}

func ensureSignedCommits(ctx context.Context, org string, repo *gh_pb.Repository, branch *gh_pb.Branch, dry bool) error {
	clt, err := client.ClientFromContext(ctx)
	if err != nil {
		return err
	}

	if branch.Protection.SignedCommits == nil {
		return nil
	}

	ghpb, err := clt.GetBranchProtection(ctx, org, repo.Name, branch.Name)
	if err != nil {
		return err
	}

	if ghpb.GetRequiredSignatures().GetEnabled() != *branch.Protection.SignedCommits {
		if dry {
			report.PrintAdd("updating require signed commits to '" + fmt.Sprintf("%t", *branch.Protection.SignedCommits) + "'")
			report.Println()

			return nil
		}

		err = clt.RequireSignedCommits(ctx, org, repo.Name, branch.Name)
		if err != nil {
			return err
		}

		report.PrintAdd("updated require signed commits to '" + fmt.Sprintf("%t", *branch.Protection.SignedCommits) + "'")
		report.Println()
	} else {
		report.PrintInfo("require signed commits is '" + fmt.Sprintf("%t", *branch.Protection.SignedCommits) + "'")
		report.Println()
	}

	return nil
}
