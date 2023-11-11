package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

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
)

type Client struct {
	ghClient *github.Client
	rate     *rate.Limiter
}

func New(ctx context.Context, tkn string) (*Client, error) {
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
