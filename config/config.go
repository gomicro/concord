package config

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type ctxKey string

const (
	ctxKeyConfig ctxKey = "config"
)

var (
	ErrConfigNotFound = fmt.Errorf("config not found")
)

type File struct {
	Github Github `yaml:"github"`
}

type Github struct {
	Token string `yaml:"token"`
}

func ParseFromFile() (*File, error) {
	file, err := GetConfigFile()
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	return parseFromFile(file)
}

func parseFromFile(file string) (*File, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDONLY, configFileMask)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	// nothing was written, so a close error can't lose data
	defer func() { _ = f.Close() }()

	var c File

	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return &c, nil
}

func (c *File) WriteToFile() error {
	file, err := GetConfigFile()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return c.writeToFile(file)
}

func (c *File) writeToFile(file string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, configFileMask)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	err = yaml.NewEncoder(f).Encode(c)
	if err != nil {
		_ = f.Close()
		return fmt.Errorf("encode: %w", err)
	}

	// close can report a failed write, so it needs to be checked
	err = f.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}

func WithConfig(ctx context.Context, file string) context.Context {
	c, err := ParseFromFile()
	if err != nil {
		ctx, cancel := context.WithCancelCause(ctx)
		cancel(err)

		return ctx
	}

	return context.WithValue(ctx, ctxKeyConfig, c)
}

func ConfigFromContext(ctx context.Context) (*File, error) {
	c, ok := ctx.Value(ctxKeyConfig).(*File)
	if !ok {
		return nil, ErrConfigNotFound
	}

	return c, nil
}
