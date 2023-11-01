package client

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/go-github/v56/github"
)

var (
	ErrOrgNotFound = errors.New("organization not found")
)

func (c *Client) GetOrg(ctx context.Context, orgName string) (*github.Organization, error) {
	org, _, err := c.ghClient.Organizations.Get(ctx, orgName)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, err
		}

		if errResp, ok := err.(*github.ErrorResponse); ok {
			if errResp.Response.StatusCode == http.StatusNotFound {
				return nil, ErrOrgNotFound
			}
		}

		return nil, err
	}

	return org, nil
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

func (c *Client) CreateTeam(ctx context.Context, orgName, teamName string) (*github.Team, error) {
	team, _, err := c.ghClient.Teams.CreateTeam(ctx, orgName, github.NewTeam{
		Name: teamName,
	})
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, err
		}

		return nil, err
	}

	return team, nil
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

func (c *Client) GetTeamMembers(ctx context.Context, orgID, teamID int64) ([]*github.User, error) {
	if teamID == -1 {
		return []*github.User{}, nil
	}

	members, _, err := c.ghClient.Teams.ListTeamMembersByID(ctx, orgID, teamID, nil)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, err
		}

		return nil, err
	}

	return members, nil
}
