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
	cs := &report.ChangeSet{}
	cs.Add("removing repo from team '"+team+"'", "removed repo from team '"+team+"'")

	cs.PrintPre()

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

		cs.PrintPost()

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
	cs := &report.ChangeSet{}
	cs.Add("creating repo "+repo.GetName(), "created repo "+repo.GetName())

	if repo.Description != nil {
		cs.Add("setting description to '"+repo.GetDescription()+"'", "set description to '"+repo.GetDescription()+"'")
	}

	if repo.Archived != nil {
		cs.Add("setting archived to '"+fmt.Sprintf("%t", repo.GetArchived())+"'", "set archived to '"+fmt.Sprintf("%t", repo.GetArchived())+"'")
	}

	if repo.Private != nil {
		cs.Add("setting private to '"+fmt.Sprintf("%t", repo.GetPrivate())+"'", "set private to '"+fmt.Sprintf("%t", repo.GetPrivate())+"'")
	}

	if repo.DefaultBranch != nil {
		cs.Add("setting default branch to '"+repo.GetDefaultBranch()+"'", "set default branch to '"+repo.GetDefaultBranch()+"'")
	}

	cs.PrintPre()

	c.Add(func() error {
		c.rate.Wait(ctx) //nolint: errcheck
		_, _, err := c.ghClient.Repositories.Create(ctx, org, repo)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return fmt.Errorf("github: hit rate limit")
			}

			return fmt.Errorf("create repo: %w", err)
		}

		cs.PrintPost()

		return nil
	})
}

func (c *Client) UpdateRepo(ctx context.Context, org, repo string, edits *github.Repository) {
	cs := &report.ChangeSet{}

	if edits.Description != nil {
		cs.Add("updating description to '"+*edits.Description+"'", "updated description to '"+*edits.Description+"'")
	}

	if edits.Archived != nil {
		cs.Add("updating archived to '"+fmt.Sprintf("%t", *edits.Archived)+"'", "updated archived to '"+fmt.Sprintf("%t", *edits.Archived)+"'")
	}

	if edits.Private != nil {
		cs.Add("updating private to '"+fmt.Sprintf("%t", *edits.Private)+"'", "updated private to '"+fmt.Sprintf("%t", *edits.Private)+"'")
	}

	if edits.DefaultBranch != nil {
		cs.Add("updating default branch to '"+*edits.DefaultBranch+"'", "updated default branch to '"+*edits.DefaultBranch+"'")
	}

	if edits.DeleteBranchOnMerge != nil {
		cs.Add("updating auto delete head branches to '"+fmt.Sprintf("%t", *edits.DeleteBranchOnMerge)+"'", "updated auto delete head branches to '"+fmt.Sprintf("%t", *edits.DeleteBranchOnMerge)+"'")
	}

	if edits.AllowAutoMerge != nil {
		cs.Add("updating allow auto merge to '"+fmt.Sprintf("%t", *edits.AllowAutoMerge)+"'", "updated allow auto merge to '"+fmt.Sprintf("%t", *edits.AllowAutoMerge)+"'")
	}

	cs.PrintPre()

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

		cs.PrintPost()

		return nil
	})
}

func (c *Client) SetRepoTopics(ctx context.Context, org, repo string, topics []string) {
	cs := &report.ChangeSet{}
	cs.Add("updating labels to ["+strings.Join(topics, ", ")+"]", "updated labels to ["+strings.Join(topics, ", ")+"]")

	cs.PrintPre()

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

		cs.PrintPost()

		return nil
	})
}

func (c *Client) ProtectBranch(ctx context.Context, org, repo, branch string, protection *github.ProtectionRequest) error {
	ghpb, err := c.GetBranchProtection(ctx, org, repo, branch)
	if err != nil && !errors.Is(err, ErrBranchProtectionNotFound) {
		return err
	}

	cs := &report.ChangeSet{}

	if ghpb != nil {
		report.PrintInfo(branch + " branch protected")
		report.Println()
	} else {
		cs.Add("protecting branch "+branch, "protected branch "+branch)
	}

	if protection.RequiredPullRequestReviews != nil {
		if ghpb.GetRequiredPullRequestReviews() == nil {
			cs.Add("setting require pr to 'true'", "set require pr to 'true'")
		}
	} else {
		if ghpb.GetRequiredPullRequestReviews() != nil {
			cs.Add("setting require pr to 'false'", "set require pr to 'false'")
		}
	}

	checks := []string{}
	if protection.RequiredStatusChecks != nil {
		if ghpb.GetRequiredStatusChecks() == nil {
			cs.Add("setting require status checks to 'true'", "set require status checks to 'true'")

			rc := protection.GetRequiredStatusChecks()
			if len(rc.Checks) > 0 {
				for i := range rc.Checks {
					checks = append(checks, rc.Checks[i].Context)
				}
			}

			if len(checks) > 0 {
				cs.Add("setting required checks to ["+strings.Join(checks, ", ")+"]", "set required checks to ["+strings.Join(checks, ", ")+"]")
			}
		} else {
			report.PrintInfo("status checks required")
			report.Println()
		}
	} else {
		if ghpb.GetRequiredStatusChecks() != nil {
			cs.Add("setting require status checks to 'false'", "set require status checks to 'false'")
		}
	}

	cs.PrintPre()

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

		cs.PrintPost()

		return nil
	})

	return nil
}

func (c *Client) SetRequireSignedCommits(ctx context.Context, org, repo, branch string, require bool) error {
	ghpb, err := c.GetBranchProtection(ctx, org, repo, branch)
	if err != nil && !errors.Is(err, ErrBranchProtectionNotFound) {
		return err
	}

	cs := &report.ChangeSet{}

	if ghpb.GetRequiredSignatures().GetEnabled() != require {
		cs.Add(fmt.Sprintf("setting require signed commits to '%t'", require), fmt.Sprintf("set require signed commits to '%t'", require))
	} else {
		report.PrintInfo(fmt.Sprintf("require signed commits is '%t'", require))
		report.Println()
	}

	cs.PrintPre()

	c.Add(func() error {
		c.rate.Wait(ctx) //nolint: errcheck
		var resp *github.Response
		var err error
		if require {
			_, resp, err = c.ghClient.Repositories.RequireSignaturesOnProtectedBranch(ctx, org, repo, branch)
		} else {
			resp, err = c.ghClient.Repositories.OptionalSignaturesOnProtectedBranch(ctx, org, repo, branch)
		}

		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return fmt.Errorf("github: hit rate limit")
			}

			if resp.StatusCode == http.StatusNotFound {
				return ErrBranchProtectionNotFound
			}

			return fmt.Errorf("protect branch: set signature required: %w", err)
		}

		cs.PrintPost()

		return nil
	})

	return nil
}
