package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

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

func (c *Client) CreateRepo(ctx context.Context, org string, repo *github.Repository) error {
	c.rate.Wait(ctx) //nolint: errcheck
	_, _, err := c.ghClient.Repositories.Create(ctx, org, repo)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return fmt.Errorf("github: hit rate limit")
		}

		return fmt.Errorf("create repo: %w", err)
	}

	return nil
}

func (c *Client) UpdateRepo(ctx context.Context, org, repo string, edits *github.Repository) error {
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

	return nil
}

func (c *Client) SetRepoTopics(ctx context.Context, org, repo string, topics []string) error {
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

	return nil
}

func (c *Client) ProtectBranch(ctx context.Context, org, repo, branch string, protection *github.ProtectionRequest) error {
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

	return nil
}

func (c *Client) RequireSignedCommits(ctx context.Context, org, repo, branch string) error {
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

	return nil
}
