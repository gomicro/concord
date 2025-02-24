package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gomicro/scribe"
	"github.com/gomicro/scribe/color"
	"github.com/gomicro/trust"
	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

const (
	BurstLimit        = 10
	RequestsPerSecond = 10
)

var (
	ErrClientNotFound = errors.New("client not found in context")
	ErrTokenEmpty     = errors.New("token is empty; please run `concord auth` or set the GITHUB_TOKEN environment variable")
)

type Client struct {
	ghClient *github.Client
	rate     *rate.Limiter

	stack []func() error
}

func New(ctx context.Context, tkn string) (*Client, error) {
	if tkn == "" {
		return nil, ErrTokenEmpty
	}

	pool := trust.New()

	certs, err := pool.CACerts()
	if err != nil {
		return nil, fmt.Errorf("failed to create cert pool: %w", err)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certs},
		},
	}

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

func (c *Client) Add(fn func() error) {
	c.stack = append(c.stack, fn)
}

func (c *Client) Apply() error {
	if len(c.stack) == 0 {
		return nil
	}

	t := &scribe.Theme{
		Describe: func(s string) string {
			return color.CyanFg(s)
		},
		Print: scribe.NoopDecorator,
	}

	scrb := scribe.NewScribe(os.Stdout, t)

	scrb.BeginDescribe("Applying")
	scrb.EndDescribe()

	for _, fn := range c.stack {
		err := fn()
		if err != nil {
			scrb.Print(color.RedFg(err.Error()))
		}
	}

	return nil
}
