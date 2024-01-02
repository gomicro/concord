package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gomicro/concord/report"
	"github.com/google/go-github/v56/github"
)

var (
	ErrOrgNotFound  = errors.New("organization not found")
	ErrUserNotFound = errors.New("user not found")
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

func (c *Client) OrgExists(ctx context.Context, orgName string) (bool, error) {
	_, err := c.GetOrg(ctx, orgName)
	if err != nil {
		if errors.Is(err, ErrOrgNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
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

func (c *Client) InviteMember(ctx context.Context, orgName string, username string) {
	cs := &report.ChangeSet{}

	cs.Add("invite "+username, "invited "+username)
	cs.PrintPre()

	c.Add(func() error {
		user, resp, err := c.ghClient.Users.Get(ctx, username)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return err
			}

			if resp.StatusCode == http.StatusNotFound {
				return ErrUserNotFound
			}

			return err
		}

		_, _, err = c.ghClient.Organizations.CreateOrgInvitation(ctx, orgName, &github.CreateOrgInvitationOptions{
			InviteeID: user.ID,
		})
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return err
			}

			return err
		}

		cs.PrintPost()

		return nil
	})
}

func (c *Client) SetOrgPrivileges(ctx context.Context, orgName string, edits *github.Organization) error {
	ghOrg, _, err := c.ghClient.Organizations.Get(ctx, orgName)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return err
		}

		if errResp, ok := err.(*github.ErrorResponse); ok {
			if errResp.Response.StatusCode == http.StatusNotFound {
				return ErrOrgNotFound
			}
		}

		return err
	}

	cs := &report.ChangeSet{}

	if edits.DefaultRepoPermission != nil && *edits.DefaultRepoPermission != *ghOrg.DefaultRepoPermission {
		cs.Add(
			fmt.Sprintf("setting base permissions to '%s'", *edits.DefaultRepoPermission),
			fmt.Sprintf("set base permissions to '%s'", *edits.DefaultRepoPermission),
		)
	}

	if edits.MembersCanCreatePrivateRepos != nil && *edits.MembersCanCreatePrivateRepos != *ghOrg.MembersCanCreatePrivateRepos {
		cs.Add(
			fmt.Sprintf("setting private repo creation to '%t'", *edits.MembersCanCreatePrivateRepos),
			fmt.Sprintf("set private repo creation to '%t'", *edits.MembersCanCreatePrivateRepos),
		)
	}

	if edits.MembersCanCreatePublicRepos != nil && *edits.MembersCanCreatePublicRepos != *ghOrg.MembersCanCreatePublicRepos {
		cs.Add(
			fmt.Sprintf("setting public repo creation to '%t'", *edits.MembersCanCreatePublicRepos),
			fmt.Sprintf("set public repo creation to '%t'", *edits.MembersCanCreatePublicRepos),
		)
	}

	cs.PrintPre()

	c.Add(func() error {
		_, resp, err := c.ghClient.Organizations.Edit(ctx, orgName, edits)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				return err
			}

			if resp.StatusCode == http.StatusNotFound {
				return ErrUserNotFound
			}

			return err
		}

		cs.PrintPost()

		return nil
	})

	return nil
}
