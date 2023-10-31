package client

import (
	"context"
	"net/http"

	"github.com/google/go-github/v56/github"
)

func (c *Client) OrgExists(ctx context.Context, orgName string) (bool, error) {
	_, resp, err := c.ghClient.Organizations.Get(ctx, orgName)
	if resp == nil && err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return false, err
		}

		return false, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

func (c *Client) GetTeams(ctx context.Context, orgName string) ([]*github.Team, error) {
	teams, _, err := c.ghClient.Teams.ListTeams(ctx, orgName, nil)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, err
		}

		return nil, err
	}

	return teams, nil
}

func (c *Client) GetMembers(ctx context.Context, orgName string) ([]*github.User, error) {
	members, _, err := c.ghClient.Organizations.ListMembers(ctx, orgName, nil)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, err
		}

		return nil, err
	}

	return members, nil
}

func (c *Client) GetTeamMembers(ctx context.Context, teamID int64) ([]*github.User, error) {
	if teamID == -1 {
		return []*github.User{}, nil
	}

	members, _, err := c.ghClient.Teams.ListTeamMembers(ctx, teamID, nil)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, err
		}

		return nil, err
	}

	return members, nil
}
