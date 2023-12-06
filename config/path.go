package config

import (
	"errors"
	"fmt"
	"os"
	"path"
)

const (
	defaultConfigFile    = "config.yml"
	defaultConfigBaseDir = ".config/concord/"

	configDirMask  = 0700
	configFileMask = 0600
)

var (
	ErrTooPermissive = errors.New("permissions too permissive")
)

func GetConfigFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("config: home: %w", err)
	}

	dir := path.Join(home, defaultConfigBaseDir)

	info, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("config: dir: %w", err)
		}

		err = os.MkdirAll(dir, configDirMask)
		if err != nil {
			return "", fmt.Errorf("config: mkdir: %w", err)
		}

		info, _ = os.Stat(dir)
	}

	if info.Mode().Perm() != configDirMask {
		return "", fmt.Errorf("config: dir mask: %w", ErrTooPermissive)
	}

	f := path.Join(dir, defaultConfigFile)

	info, err = os.Stat(f)
	if err != nil && !os.IsNotExist(err) {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("config: file: %w", err)
		}
	}

	if info != nil && info.Mode().Perm() != configFileMask {
		return "", fmt.Errorf("config: file mask: %w", ErrTooPermissive)
	}

	return f, nil
}
