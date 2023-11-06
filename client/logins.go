package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v56/github"
)

func (c *Client) GetLogins(ctx context.Context) ([]string, error) {
	logins := []string{}

	user, _, err := c.ghClient.Users.Get(ctx, "")
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("github: hit rate limit")
		}

		return nil, fmt.Errorf("get user: %w", err)
	}

	logins = append(logins, strings.ToLower(user.GetLogin()))

	opts := &github.ListOptions{
		Page:    0,
		PerPage: 100,
	}

	orgs, _, err := c.ghClient.Organizations.List(ctx, "", opts)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("github: hit rate limit")
		}

		return nil, fmt.Errorf("list orgs: %w", err)
	}

	for i := range orgs {
		o := orgs[i].GetLogin()
		logins = append(logins, strings.ToLower(o))
	}

	return logins, nil
}
