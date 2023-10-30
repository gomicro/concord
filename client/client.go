package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/gomicro/trust"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

const (
	BurstLimit        = 10
	RequestsPerSecond = 10
)

type Client struct {
	ghClient *github.Client
	rate     *rate.Limiter
}

func New(tkn string) (*Client, error) {
	pool := trust.New()

	certs, err := pool.CACerts()
	if err != nil {
		return nil, fmt.Errorf("failed to create cert pool: %v\n", err.Error())
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certs},
		},
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: tkn,
		},
	)

	rl := rate.NewLimiter(
		rate.Limit(RequestsPerSecond),
		BurstLimit,
	)

	return &Client{
		ghClient: github.NewClient(oauth2.NewClient(ctx, ts)),
		rate:     rl,
	}, nil
}

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