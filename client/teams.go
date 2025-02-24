package client

import (
	"context"

	"github.com/gomicro/scribe"
	"github.com/gomicro/scribe/color"
	"github.com/google/go-github/v56/github"
)

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

func (c *Client) CreateTeam(ctx context.Context, scrb scribe.Scriber, orgName, teamName string) {
	scrb.Print(color.GreenFg("create team " + teamName))

	c.Add(func() error {
		team, _, err := c.ghClient.Teams.CreateTeam(ctx, orgName, github.NewTeam{
			Name: teamName,
		})
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return err
			}

			return err
		}

		// when creating a team, the current user is added, so we need to remove it
		user, _, err := c.ghClient.Users.Get(ctx, "")
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return err
			}

			return err
		}

		err = c.RemoveTeamMember(ctx, team.GetOrganization().GetID(), team.GetID(), *user.Login)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return err
			}

			return err
		}

		scrb.Print(color.GreenFg("created team " + teamName))

		return nil
	})
}

func (c *Client) InviteTeamMember(ctx context.Context, scrb scribe.Scriber, org, team, user string) {
	scrb.Print(color.GreenFg("invite " + user + " to team " + team))

	c.Add(func() error {
		_, _, err := c.ghClient.Teams.AddTeamMembershipBySlug(ctx, org, team, user, nil)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return err
			}

			return err
		}

		scrb.Print(color.GreenFg("invited " + user + " to team " + team))

		return nil
	})
}

func (c *Client) RemoveTeamMember(ctx context.Context, orgID, teamID int64, user string) error {
	_, err := c.ghClient.Teams.RemoveTeamMembershipByID(ctx, orgID, teamID, user)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return err
		}

		return err
	}

	return nil
}

func (c *Client) GetTeamMembers(ctx context.Context, org, team string) ([]*github.User, error) {
	members, _, err := c.ghClient.Teams.ListTeamMembersBySlug(ctx, org, team, nil)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, err
		}

		return nil, err
	}

	return members, nil
}
