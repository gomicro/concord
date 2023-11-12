package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gomicro/concord/report"
	"github.com/google/go-github/v56/github"
)

var (
	ErrGetBranch                = errors.New("get branch")
	ErrRepoNotFound             = errors.New("repo not found")
	ErrNoReposFound             = errors.New("no repos found")
	ErrBranchProtectionNotFound = errors.New("branch protection not found")
)

func (c *Client) GetRepos(ctx context.Context, name string) ([]*github.Repository, error) {
	count := int64(0)
	orgFound := true

	c.rate.Wait(ctx) //nolint: errcheck
	org, resp, err := c.ghClient.Organizations.Get(ctx, name)
	if resp == nil && err != nil {

		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("github: hit rate limit")
		}

		return nil, fmt.Errorf("get org: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		orgFound = false

		c.rate.Wait(ctx) //nolint: errcheck
		user, _, err := c.ghClient.Users.Get(ctx, name)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return nil, fmt.Errorf("github: hit rate limit")
			}

			return nil, fmt.Errorf("get user: %v", err.Error())
		}

		count = int64(user.GetPublicRepos()) + user.GetTotalPrivateRepos()
	} else {
		count = int64(org.GetPublicRepos()) + org.GetTotalPrivateRepos()
	}

	if count < 1 {
		return nil, ErrNoReposFound
	}

	orgOpts := &github.RepositoryListByOrgOptions{
		Type: "all",
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}

	userOpts := &github.RepositoryListOptions{
		Type: "all",
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}

	var repos []*github.Repository
	for {
		var rs []*github.Repository
		c.rate.Wait(ctx) //nolint: errcheck
		if orgFound {
			rs, resp, err = c.ghClient.Repositories.ListByOrg(ctx, name, orgOpts)
		} else {
			rs, resp, err = c.ghClient.Repositories.List(ctx, name, userOpts)
		}

		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return nil, fmt.Errorf("github: hit rate limit")
			}

			return nil, fmt.Errorf("list repos: %v", err.Error())
		}

		for i := range rs {
			if rs[i].GetArchived() {
				continue
			}

			repos = append(repos, rs[i])
		}

		if resp.NextPage == 0 {
			break
		}

		if orgFound {
			orgOpts.Page = resp.NextPage
		} else {
			userOpts.Page = resp.NextPage
		}
	}

	return repos, nil
}

func (c *Client) GetRepo(ctx context.Context, org, name string) (*github.Repository, error) {
	c.rate.Wait(ctx) //nolint: errcheck
	repo, resp, err := c.ghClient.Repositories.Get(ctx, org, name)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("github: hit rate limit")
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrRepoNotFound
		}

		return nil, fmt.Errorf("get repo: %w", err)
	}

	return repo, nil
}

func (c *Client) GetRepoTeams(ctx context.Context, org, repo string) ([]*github.Team, error) {
	c.rate.Wait(ctx) //nolint: errcheck
	teams, resp, err := c.ghClient.Repositories.ListTeams(ctx, org, repo, nil)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("github: hit rate limit")
		}

		return nil, fmt.Errorf("get repo teams: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrRepoNotFound
	}

	return teams, nil
}

func (c *Client) AddRepoToTeam(ctx context.Context, org, team, repo, perm string) error {
	gts, err := c.GetRepoTeams(ctx, org, repo)
	if err != nil {
		return err
	}

	gtps := map[string]string{}
	for _, gt := range gts {
		gtps[strings.ToLower(gt.GetName())] = gt.GetPermission()
	}

	tp, relationExists := gtps[strings.ToLower(team)]

	p := perm
	switch p {
	case "read":
		p = "pull"
	case "write":
		p = "push"
	}

	if relationExists && strings.EqualFold(tp, p) {
		report.PrintInfo("team '" + team + "' has permission '" + perm + "'")
		report.Println()
		return nil
	}

	report.PrintAdd("adding repo to team '" + team + "' with '" + perm + "'")
	report.Println()

	c.Add(func() error {
		c.rate.Wait(ctx) //nolint: errcheck

		resp, err := c.ghClient.Teams.AddTeamRepoBySlug(ctx, org, team, org, repo, &github.TeamAddTeamRepoOptions{
			Permission: p,
		})
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return fmt.Errorf("github: hit rate limit")
			}

			if resp.StatusCode == http.StatusNotFound {
				return ErrRepoNotFound
			}

			return fmt.Errorf("add repo to team: %w", err)
		}

		if relationExists {
			report.PrintSuccess("updated repo to team '" + team + "' with '" + perm + "'")
		} else {
			report.PrintAdd("added repo to team '" + team + "' with '" + perm + "'")
			report.Println()
		}

		return nil
	})

	return nil
}

func (c *Client) RemoveRepoFromTeam(ctx context.Context, org, team, repo string) {
	report.PrintAdd("removing repo from team '" + team + "'")
	report.Println()

	c.Add(func() error {
		c.rate.Wait(ctx) //nolint: errcheck
		resp, err := c.ghClient.Teams.RemoveTeamRepoBySlug(ctx, org, team, org, repo)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return fmt.Errorf("github: hit rate limit")
			}
			if resp.StatusCode == http.StatusNotFound {
				return ErrRepoNotFound
			}

			return fmt.Errorf("remove repo from team: %w", err)
		}

		report.PrintAdd("removed repo from team '" + team + "'")
		report.Println()

		return nil
	})
}

func (c *Client) GetRepoTopics(ctx context.Context, org, name string) ([]string, error) {
	c.rate.Wait(ctx) //nolint: errcheck
	topics, resp, err := c.ghClient.Repositories.ListAllTopics(ctx, org, name)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("github: hit rate limit")
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrRepoNotFound
		}

		return nil, fmt.Errorf("get repo topics: %w", err)
	}

	return topics, nil
}

func (c *Client) GetBranches(ctx context.Context, org, repo string) ([]*github.Branch, error) {
	c.rate.Wait(ctx) //nolint: errcheck
	branches, resp, err := c.ghClient.Repositories.ListBranches(ctx, org, repo, nil)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("github: hit rate limit")
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrRepoNotFound
		}

		return nil, fmt.Errorf("get branches: %w", err)
	}

	return branches, nil
}

func (c *Client) GetBranchProtection(ctx context.Context, org, repo, branch string) (*github.Protection, error) {
	c.rate.Wait(ctx) //nolint: errcheck
	b, resp, err := c.ghClient.Repositories.GetBranchProtection(ctx, org, repo, branch)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("github: hit rate limit")
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrBranchProtectionNotFound
		}

		return nil, fmt.Errorf("get branch: %w", err)
	}

	return b, nil
}

func (c *Client) IsBranchProtected(ctx context.Context, org, repo, branch string) (bool, error) {
	c.rate.Wait(ctx) //nolint: errcheck
	b, resp, err := c.ghClient.Repositories.GetBranchProtection(ctx, org, repo, branch)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return false, fmt.Errorf("github: hit rate limit")
		}

		if resp.StatusCode == http.StatusNotFound {
			return false, nil
		}

		return false, fmt.Errorf("get branch: %w", err)
	}

	return b != nil, nil
}

func (c *Client) CreateRepo(ctx context.Context, org string, repo *github.Repository) {
	report.PrintAdd("creating repo " + repo.GetName())
	report.Println()

	if repo.Description != nil {
		report.PrintAdd("setting description to '" + repo.GetDescription() + "'")
		report.Println()
	}

	if repo.Archived != nil {
		report.PrintAdd("setting archived to '" + fmt.Sprintf("%t", repo.GetArchived()) + "'")
		report.Println()
	}

	if repo.Private != nil {
		report.PrintAdd("setting private to '" + fmt.Sprintf("%t", repo.GetPrivate()) + "'")
		report.Println()
	}

	if repo.DefaultBranch != nil {
		report.PrintAdd("setting default branch to '" + repo.GetDefaultBranch() + "'")
		report.Println()
	}

	c.Add(func() error {
		c.rate.Wait(ctx) //nolint: errcheck
		_, _, err := c.ghClient.Repositories.Create(ctx, org, repo)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return fmt.Errorf("github: hit rate limit")
			}

			return fmt.Errorf("create repo: %w", err)
		}

		report.PrintSuccess("created repo " + repo.GetName())
		report.Println()

		if repo.Description != nil {
			report.PrintSuccess("set description to '" + repo.GetDescription() + "'")
			report.Println()
		}

		if repo.Archived != nil {
			report.PrintSuccess("set archived to '" + fmt.Sprintf("%t", repo.GetArchived()) + "'")
			report.Println()
		}

		if repo.Private != nil {
			report.PrintSuccess("set private to '" + fmt.Sprintf("%t", repo.GetPrivate()) + "'")
			report.Println()
		}

		if repo.DefaultBranch != nil {
			report.PrintSuccess("set default branch to '" + repo.GetDefaultBranch() + "'")
			report.Println()
		}

		return nil
	})
}

func (c *Client) UpdateRepo(ctx context.Context, org, repo string, edits *github.Repository) {
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

	c.Add(func() error {
		c.rate.Wait(ctx) //nolint: errcheck
		_, resp, err := c.ghClient.Repositories.Edit(ctx, org, repo, edits)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return fmt.Errorf("github: hit rate limit")
			}

			if resp.StatusCode == http.StatusNotFound {
				return ErrRepoNotFound
			}

			return fmt.Errorf("update repo description: %w", err)
		}

		if edits.Description != nil {
			report.PrintSuccess("updated description to '" + *edits.Description + "'")
			report.Println()
		}

		if edits.Archived != nil {
			report.PrintSuccess("updated archived to '" + fmt.Sprintf("%t", *edits.Archived) + "'")
			report.Println()
		}

		if edits.Private != nil {
			report.PrintSuccess("updated private to '" + fmt.Sprintf("%t", *edits.Private) + "'")
			report.Println()
		}

		if edits.DefaultBranch != nil {
			report.PrintSuccess("updated default branch to '" + *edits.DefaultBranch + "'")
			report.Println()
		}

		if edits.DeleteBranchOnMerge != nil {
			report.PrintSuccess("updated auto delete head branches to '" + fmt.Sprintf("%t", *edits.DeleteBranchOnMerge) + "'")
			report.Println()
		}

		if edits.AllowAutoMerge != nil {
			report.PrintSuccess("updated allow auto merge to '" + fmt.Sprintf("%t", *edits.AllowAutoMerge) + "'")
			report.Println()
		}

		return nil
	})
}

func (c *Client) SetRepoTopics(ctx context.Context, org, repo string, topics []string) {
	report.PrintAdd("updating labels to [" + strings.Join(topics, ", ") + "]")
	report.Println()

	c.Add(func() error {
		c.rate.Wait(ctx) //nolint: errcheck
		_, resp, err := c.ghClient.Repositories.ReplaceAllTopics(ctx, org, repo, topics)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return fmt.Errorf("github: hit rate limit")
			}

			if resp.StatusCode == http.StatusNotFound {
				return ErrRepoNotFound
			}

			return fmt.Errorf("set repo topics: %w", err)
		}

		report.PrintSuccess("updated labels to [" + strings.Join(topics, ", ") + "]")
		report.Println()

		return nil
	})
}

func (c *Client) ProtectBranch(ctx context.Context, org, repo, branch string, protection *github.ProtectionRequest) {
	report.PrintAdd("protecting branch " + branch)
	report.Println()

	if protection.RequiredPullRequestReviews != nil {
		report.PrintAdd("setting require pr to 'true'")
		report.Println()
	}

	checks := []string{}
	if protection.RequiredStatusChecks != nil {
		report.PrintAdd("setting require status checks to 'true'")
		report.Println()

		rc := protection.GetRequiredStatusChecks()
		if len(rc.Checks) > 0 {
			for i := range rc.Checks {
				checks = append(checks, rc.Checks[i].Context)
			}
		}

		if len(checks) > 0 {
			report.PrintAdd("setting required checks to [" + strings.Join(checks, ", ") + "]")
			report.Println()
		}
	}

	c.Add(func() error {
		c.rate.Wait(ctx) //nolint: errcheck
		_, resp, err := c.ghClient.Repositories.UpdateBranchProtection(ctx, org, repo, branch, protection)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return fmt.Errorf("github: hit rate limit")
			}

			if resp.StatusCode == http.StatusNotFound {
				return ErrBranchProtectionNotFound
			}

			return fmt.Errorf("protect branch: %w", err)
		}

		report.PrintSuccess("protected branch " + branch)
		report.Println()

		if protection.RequiredPullRequestReviews != nil {
			report.PrintSuccess("set require pr to 'true'")
			report.Println()
		} else {
			report.PrintSuccess("set require pr to 'false'")
			report.Println()
		}

		if protection.RequiredStatusChecks != nil {
			report.PrintSuccess("set require status checks to 'true'")
			report.Println()

			if len(checks) > 0 {
				report.PrintSuccess("set required checks to [" + strings.Join(checks, ", ") + "]")
				report.Println()
			}
		} else {
			report.PrintSuccess("set require status checks to 'false'")
			report.Println()
		}

		return nil
	})
}

func (c *Client) RequireSignedCommits(ctx context.Context, org, repo, branch string) {
	report.PrintAdd("updating require signed commits to 'true'")
	report.Println()

	c.Add(func() error {
		c.rate.Wait(ctx) //nolint: errcheck
		_, resp, err := c.ghClient.Repositories.RequireSignaturesOnProtectedBranch(ctx, org, repo, branch)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return fmt.Errorf("github: hit rate limit")
			}

			if resp.StatusCode == http.StatusNotFound {
				return ErrBranchProtectionNotFound
			}

			return fmt.Errorf("protect branch: signature required: %w", err)
		}

		report.PrintSuccess("updated require signed commits to 'true'")
		report.Println()

		return nil
	})
}
