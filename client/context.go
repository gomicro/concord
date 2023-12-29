package client

import "context"

type ctxKey string

const (
	clientConextKey ctxKey = "client"
)

func WithClient(ctx context.Context, tkn string) (context.Context, error) {
	c, err := New(ctx, tkn)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, clientConextKey, c), nil
}

func ClientFromContext(ctx context.Context) (*Client, error) {
	c, ok := ctx.Value(clientConextKey).(*Client)
	if !ok {
		return nil, ErrClientNotFound
	}

	return c, nil
}
