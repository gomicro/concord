package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gomicro/scribe"
	"github.com/gomicro/scribe/color"
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

		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrRepoNotFound
		}

		return nil, fmt.Errorf("get repo teams: %w", err)
	}

	return teams, nil
}

func (c *Client) AddRepoToTeam(ctx context.Context, scrb scribe.Scriber, org, team, repo, perm string) error {
	gts, err := c.GetRepoTeams(ctx, org, repo)
	if err != nil {
		return fmt.Errorf("add repo team: %w", err)
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
		scrb.Print("team '" + team + "' already has permission '" + perm + "'")
		return nil
	}

	scrb.Print(color.GreenFg("adding repo to team '" + team + "' with '" + perm + "'"))

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

			return fmt.Errorf("%s/%s: add repo to team: %w", org, repo, err)
		}

		if relationExists {
			scrb.Print(color.GreenFg("updated repo to team '" + team + "' with '" + perm + "'"))
		} else {
			scrb.Print(color.GreenFg("added repo to team '" + team + "' with '" + perm + "'"))
		}

		return nil
	})

	return nil
}

func (c *Client) RemoveRepoFromTeam(ctx context.Context, scrb scribe.Scriber, org, team, repo string) {
	scrb.Print(color.GreenFg("removing repo from team '" + team + "'"))

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

			return fmt.Errorf("%s/%s: remove repo from team: %w", org, repo, err)
		}

		scrb.Print(color.GreenFg("removed repo from team '" + team + "'"))

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

func (c *Client) CreateRepo(ctx context.Context, scrb scribe.Scriber, org string, repo *github.Repository) {
	scrb.Print(color.GreenFg("creating repo " + repo.GetName()))

	if repo.Description != nil {
		scrb.Print(color.GreenFg("set description to '" + repo.GetDescription() + "'"))
	}

	if repo.Archived != nil {
		scrb.Print(color.GreenFg("set archived to '" + fmt.Sprintf("%t", repo.GetArchived()) + "'"))
	}

	if repo.Private != nil {
		scrb.Print(color.GreenFg("set private to '" + fmt.Sprintf("%t", repo.GetPrivate()) + "'"))
	}

	if repo.DefaultBranch != nil {
		scrb.Print(color.GreenFg("set default branch to '" + repo.GetDefaultBranch() + "'"))
	}

	c.Add(func() error {
		c.rate.Wait(ctx) //nolint: errcheck
		_, _, err := c.ghClient.Repositories.Create(ctx, org, repo)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return fmt.Errorf("github: hit rate limit")
			}

			return fmt.Errorf("%s/%s: create repo: %w", org, repo, err)
		}

		scrb.Print(color.GreenFg("created repo " + repo.GetName()))

		if repo.Description != nil {
			scrb.Print(color.GreenFg("set description to '" + repo.GetDescription() + "'"))
		}

		if repo.Archived != nil {
			scrb.Print(color.GreenFg("set archived to '" + fmt.Sprintf("%t", repo.GetArchived()) + "'"))
		}

		if repo.Private != nil {
			scrb.Print(color.GreenFg("set private to '" + fmt.Sprintf("%t", repo.GetPrivate()) + "'"))
		}

		if repo.DefaultBranch != nil {
			scrb.Print(color.GreenFg("set default branch to '" + repo.GetDefaultBranch() + "'"))
		}

		return nil
	})
}

func (c *Client) InitRepo(ctx context.Context, scrb scribe.Scriber, org, repo, branch string) {
	filename := "README.md"
	content := []byte("# " + repo)

	opts := &github.RepositoryContentFileOptions{
		Message: github.String("initializing repo"),
		Content: content,
		Branch:  &branch,
	}

	scrb.Print(color.GreenFg("creating file " + filename))

	c.Add(func() error {
		c.rate.Wait(ctx) //nolint: errcheck
		_, _, err := c.ghClient.Repositories.CreateFile(ctx, org, repo, filename, opts)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return fmt.Errorf("github: hit rate limit")
			}

			return fmt.Errorf("%s/%s: init repo: %w", org, repo, err)
		}

		scrb.Print(color.GreenFg("created file " + filename))

		return nil
	})
}

func (c *Client) UpdateRepo(ctx context.Context, scrb scribe.Scriber, org, repo string, edits *github.Repository) {
	changes := false
	if edits.Description != nil {
		scrb.Print(color.GreenFg("updating description to '" + *edits.Description + "'"))
		changes = true
	}

	if edits.Archived != nil {
		scrb.Print(color.GreenFg("updating archived to '" + fmt.Sprintf("%t", *edits.Archived) + "'"))
		changes = true
	}

	if edits.Private != nil {
		scrb.Print(color.GreenFg("updating private to '" + fmt.Sprintf("%t", *edits.Private) + "'"))
		changes = true
	}

	if edits.DefaultBranch != nil {
		scrb.Print(color.GreenFg("updating default branch to '" + *edits.DefaultBranch + "'"))
		changes = true
	}

	if edits.DeleteBranchOnMerge != nil {
		scrb.Print(color.GreenFg("updating auto delete head branches to '" + fmt.Sprintf("%t", *edits.DeleteBranchOnMerge) + "'"))
		changes = true
	}

	if edits.AllowAutoMerge != nil {
		scrb.Print(color.GreenFg("updating allow auto merge to '" + fmt.Sprintf("%t", *edits.AllowAutoMerge) + "'"))
		changes = true
	}

	if !changes {
		return
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

			if resp.StatusCode == http.StatusForbidden && strings.Contains(err.Error(), "was archived so is read-only") {
				return fmt.Errorf("%s/%s: update repo: repo is archived", org, repo)
			}

			return fmt.Errorf("%s/%s: update repo: %w", org, repo, err)
		}

		if edits.Description != nil {
			scrb.Print(color.GreenFg("updated description to '" + *edits.Description + "'"))
		}

		if edits.Archived != nil {
			scrb.Print(color.GreenFg("updated archived to '" + fmt.Sprintf("%t", *edits.Archived) + "'"))
		}

		if edits.Private != nil {
			scrb.Print(color.GreenFg("updated private to '" + fmt.Sprintf("%t", *edits.Private) + "'"))
		}

		if edits.DefaultBranch != nil {
			scrb.Print(color.GreenFg("updated default branch to '" + *edits.DefaultBranch + "'"))
		}

		if edits.DeleteBranchOnMerge != nil {
			scrb.Print(color.GreenFg("updated auto delete head branches to '" + fmt.Sprintf("%t", *edits.DeleteBranchOnMerge) + "'"))
		}

		if edits.AllowAutoMerge != nil {
			scrb.Print(color.GreenFg("updated allow auto merge to '" + fmt.Sprintf("%t", *edits.AllowAutoMerge) + "'"))
		}

		return nil
	})
}

func (c *Client) SetRepoTopics(ctx context.Context, scrb scribe.Scriber, org, repo string, topics []string) {
	scrb.Print(color.GreenFg("updating labels to [" + strings.Join(topics, ", ") + "]"))

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

			return fmt.Errorf("%s/%s: set repo topics: %w", org, repo, err)
		}

		scrb.Print(color.GreenFg("updated labels to [" + strings.Join(topics, ", ") + "]"))

		return nil
	})
}

func (c *Client) ProtectBranch(ctx context.Context, scrb scribe.Scriber, org, repo, branch string, protection *github.ProtectionRequest) error {
	ghpb, err := c.GetBranchProtection(ctx, org, repo, branch)
	if err != nil && !errors.Is(err, ErrBranchProtectionNotFound) {
		return err
	}

	protecting := false
	if ghpb != nil {
		scrb.Print(branch + " branch protected")
	} else {
		scrb.Print(color.GreenFg("protecting branch " + branch))
	}

	setReqPR := false
	setReqPRTo := false
	if protection.RequiredPullRequestReviews != nil {
		if ghpb.GetRequiredPullRequestReviews() == nil {
			scrb.Print(color.GreenFg("setting require pr to 'true'"))
			setReqPR = true
			setReqPRTo = true
		}
	} else {
		if ghpb.GetRequiredPullRequestReviews() != nil {
			scrb.Print(color.GreenFg("setting require pr to 'false'"))
			setReqPR = true
			setReqPRTo = false
		}
	}

	setReqChecks := false
	setReqChecksTo := false
	setChecks := false

	checks := []string{}
	if protection.RequiredStatusChecks != nil {
		if ghpb.GetRequiredStatusChecks() == nil {
			scrb.Print(color.GreenFg("setting require status checks to 'true'"))
			setReqChecks = true
			setReqChecksTo = true

			rc := protection.GetRequiredStatusChecks()
			if len(rc.Checks) > 0 {
				for i := range rc.Checks {
					checks = append(checks, rc.Checks[i].Context)
				}
			}

			if len(checks) > 0 {
				scrb.Print(color.GreenFg("setting required checks to [" + strings.Join(checks, ", ") + "]"))
				setChecks = true
			}
		} else {
			scrb.Print("status checks required")
		}
	} else {
		if ghpb.GetRequiredStatusChecks() != nil {
			scrb.Print(color.GreenFg("setting require status checks to 'false'"))
			setReqChecks = true
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

			if resp.StatusCode == http.StatusForbidden && strings.Contains(err.Error(), "was archived so is read-only") {
				return fmt.Errorf("%s/%s: protect branch: repo is archived", org, repo)
			}

			return fmt.Errorf("%s/%s: protect branch: %w", org, repo, err)
		}

		if protecting {
			scrb.Print(color.GreenFg("protected branch " + branch))
		}

		if setReqPR {
			if setReqPRTo {
				scrb.Print(color.GreenFg("set require pr to 'true'"))
			} else {
				scrb.Print(color.GreenFg("set require pr to 'false'"))
			}
		}

		if setReqChecks {
			if setReqChecksTo {
				scrb.Print(color.GreenFg("set require status checks to 'true'"))
			} else {
				scrb.Print(color.GreenFg("set require status checks to 'false'"))
			}
		}

		if setChecks {
			scrb.Print(color.GreenFg("set required checks to [" + strings.Join(checks, ", ") + "]"))
		}

		return nil
	})

	return nil
}

func (c *Client) SetRequireSignedCommits(ctx context.Context, scrb scribe.Scriber, org, repo, branch string, require bool) error {
	ghpb, err := c.GetBranchProtection(ctx, org, repo, branch)
	if err != nil && !errors.Is(err, ErrBranchProtectionNotFound) {
		return err
	}

	req := false
	if ghpb.GetRequiredSignatures().GetEnabled() != require {
		scrb.Print(color.GreenFg(fmt.Sprintf("setting require signed commits to '%t'", require)))
		req = true
	} else {
		scrb.Print(fmt.Sprintf("require signed commits is '%t'", require))
	}

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

			if resp.StatusCode == http.StatusForbidden && strings.Contains(err.Error(), "was archived so is read-only") {
				return fmt.Errorf("%s/%s: protect branch: set signature required: repo is archived", org, repo)
			}

			return fmt.Errorf("%s/%s: protect branch: set signature required: %w", org, repo, err)
		}

		if req {
			scrb.Print(color.GreenFg(fmt.Sprintf("set require signed commits to '%t'", require)))
		}

		return nil
	})

	return nil
}
