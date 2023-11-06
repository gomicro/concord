package cmd

import (
	"context"
	"os"

	"github.com/gomicro/concord/client"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initEnvs)
}

func initEnvs() {
}

var rootCmd = &cobra.Command{
	Use:   "concord",
	Short: "concord is a tool to manage your Github repositories",
}

func Execute() {
	// TODO: Read token from config file or env, with env taking precedence
	tkn := os.Getenv("GITHUB_TOKEN")

	ctx := client.WithClient(context.Background(), tkn)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func handleError(c *cobra.Command, err error) error {
	c.SilenceUsage = true
	return err
}
