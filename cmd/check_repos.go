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

func checkRepos(ctx context.Context, manifestRepos []*gh_pb.Repository, githubRepos []*github.Repository) []*gh_pb.Repository {
	missing := []*gh_pb.Repository{}

	for _, mr := range manifestRepos {
		found := false
		for _, gr := range githubRepos {
			if strings.EqualFold(mr.Name, *gr.Name) {
				found = true
				break
			}
		}

		if !found {
			missing = append(missing, mr)
		}
	}

	return missing
}

func createRepo(ctx context.Context, org string, repo *gh_pb.Repository, dry bool) error {
	if dry {
		report.PrintWarn("create repo " + repo.Name)
		report.Println()
	}

	/*
		_, err := clt.CreateRepo(ctx, org, r)
		if err != nil {
			return err
		}
	*/

	return nil
}

func ensureRepo(ctx context.Context, org string, repo *gh_pb.Repository, dry bool) error {
	r, err := clt.GetRepo(ctx, org, repo.Name)
	if err != nil && !errors.Is(err, client.ErrRepoNotFound) {
		return err
	}

	creating := false
	if errors.Is(err, client.ErrRepoNotFound) {
		creating = true
		err = createRepo(ctx, org, repo, dry)
		if err != nil {
			return err
		}
	}

	// description
	err = ensureDescription(ctx, repo, r, creating, dry)
	if err != nil {
		return err
	}

	// archived
	err = ensureArchived(ctx, repo, r, creating, dry)
	if err != nil {
		return err
	}

	// labels
	err = ensureLabels(ctx, org, repo, r, creating, dry)
	if err != nil {
		return err
	}

	// files
	err = ensureFiles(ctx, org, repo, r, creating, dry)
	if err != nil {
		return err
	}

	// private
	err = ensurePrivate(ctx, repo, r, creating, dry)
	if err != nil {
		return err
	}

	// default branch
	err = ensureDefaultBranch(ctx, repo, r, creating, dry)
	if err != nil {
		return err
	}

	// protected branches
	err = ensureProtectedBranches(ctx, org, repo, r, creating, dry)
	if err != nil {
		return err
	}

	return nil
}

func ensureDescription(ctx context.Context, repo *gh_pb.Repository, r *github.Repository, creating, dry bool) error {
	if repo.Description == "" {
		return nil
	}

	if creating && dry {
		report.PrintWarn("update description for repo " + repo.Name)
		report.Println()
		return nil
	}

	if !strings.EqualFold(repo.Description, r.GetDescription()) {
		if dry {
			report.PrintWarn("update description for repo " + repo.Name)
			report.Println()
			return nil
		}

		/*
			_, _, err := clt.UpdateRepo(ctx, repo.Name, &github.Repository{
				Description: &repo.Description,
			})
			if err != nil {
				return handleError(cmd, err)
			}
		*/
	}

	return nil
}

func ensureArchived(ctx context.Context, repo *gh_pb.Repository, r *github.Repository, creating, dry bool) error {
	if repo.Archived == nil {
		return nil
	}

	if creating && dry {
		report.PrintWarn("archive repo " + repo.Name)
		report.Println()
		return nil
	}

	if *repo.Archived {
		if !r.GetArchived() {
			if dry {
				report.PrintWarn("archive repo " + repo.Name)
				report.Println()
				return nil
			}

			/*
				_, err := clt.ArchiveRepo(ctx, repo.Name)
				if err != nil {
					return handleError(cmd, err)
				}
			*/
		}
	} else {
		if r.GetArchived() {
			if dry {
				report.PrintWarn("unarchive repo " + repo.Name)
				report.Println()
				return nil
			}

			/*
				_, err := clt.UnarchiveRepo(ctx, repo.Name)
				if err != nil {
					return handleError(cmd, err)
				}
			*/
		}
	}

	return nil
}

func ensureFiles(ctx context.Context, org string, repo *gh_pb.Repository, r *github.Repository, creating, dry bool) error {
	return nil
}

func ensureLabels(ctx context.Context, org string, repo *gh_pb.Repository, r *github.Repository, creating, dry bool) error {
	if len(repo.Labels) == 0 {
		return nil
	}

	if creating && dry {
		for _, label := range repo.Labels {
			report.PrintWarn("create label " + label + " for repo " + repo.Name)
			report.Println()
			continue
		}

		return nil
	}

	// get existing labels
	topics, err := clt.GetRepoTopics(ctx, org, repo.Name)
	if err != nil {
		return err
	}

	// check for missing labels
	for _, label := range repo.Labels {
		found := false
		for _, topic := range topics {
			if strings.EqualFold(label, topic) {
				found = true
				break
			}
		}

		if !found {
			if dry {
				report.PrintWarn("create label " + label + " for repo " + repo.Name)
				report.Println()
				continue
			}

			/*
				_, _, err := clt.CreateLabel(ctx, repo.Name, &github.Label{
					Name:  &l.Name,
					Color: &l.Color,
				})
				if err != nil {
					return handleError(cmd, err)
				}
			*/
		}
	}

	// check for extra labels
	for _, topic := range topics {
		found := false
		for _, label := range repo.Labels {
			if strings.EqualFold(label, topic) {
				found = true
				break
			}
		}

		if !found {
			if dry {
				report.PrintWarn("delete label " + topic + " for repo " + repo.Name)
				report.Println()
				continue
			}

			/*
				_, err := clt.DeleteLabel(ctx, repo.Name, *el.Name)
				if err != nil {
					return handleError(cmd, err)
				}
			*/
		}
	}

	return nil
}

func ensurePrivate(ctx context.Context, repo *gh_pb.Repository, r *github.Repository, creating, dry bool) error {
	if repo.Private == nil {
		return nil
	}
	if creating && dry {
		report.PrintWarn("make repo " + repo.Name + " private")
		report.Println()
		return nil
	}

	if *repo.Private {
		if !r.GetPrivate() {
			if dry {
				report.PrintWarn("make repo " + repo.Name + " private")
				report.Println()
				return nil
			}

			/*
				_, _, err := clt.UpdateRepo(ctx, repo.Name, &github.Repository{
					Private: &repo.Private,
				})
				if err != nil {
					return handleError(cmd, err)
				}
			*/
		}
	} else {
		if r.GetPrivate() {
			if dry {
				report.PrintWarn("make repo " + repo.Name + " public")
				report.Println()
				return nil
			}

			/*
				_, _, err := clt.UpdateRepo(ctx, repo.Name, &github.Repository{
					Private: &repo.Private,
				})
				if err != nil {
					return handleError(cmd, err)
				}
			*/
		}
	}

	return nil
}

func ensureDefaultBranch(ctx context.Context, repo *gh_pb.Repository, r *github.Repository, creating, dry bool) error {
	if creating && dry {
		report.PrintWarn("update default branch for repo " + repo.Name)
		report.Println()
		return nil
	}

	if !strings.EqualFold(*repo.DefaultBranch, r.GetDefaultBranch()) {
		if dry {
			report.PrintWarn("update default branch for repo " + repo.Name)
			report.Println()
			return nil
		}

		/*
			_, _, err := clt.UpdateRepo(ctx, repo.Name, &github.Repository{
				DefaultBranch: &repo.DefaultBranch,
			})
			if err != nil {
				return handleError(cmd, err)
			}
		*/
	}

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
