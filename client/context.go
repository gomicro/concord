package client

import "context"

type ctxKey string

const (
	clientConextKey ctxKey = "client"
)

func WithClient(ctx context.Context, tkn string) context.Context {
	ctx, cancel := context.WithCancelCause(ctx)

	c, err := New(ctx, tkn)
	if err != nil {
		cancel(err)
	}

	return context.WithValue(ctx, clientConextKey, c)
}

func ClientFromContext(ctx context.Context) (*Client, error) {
	c, ok := ctx.Value(clientConextKey).(*Client)
	if !ok {
		return nil, ErrClientNotFound
	}

	return c, nil
}
