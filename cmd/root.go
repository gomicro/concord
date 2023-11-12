package cmd

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/gomicro/concord/client"
	"github.com/gomicro/concord/report"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initEnvs)

	rootCmd.PersistentFlags().StringP("file", "f", "concord.yml", "Path to a file containing a manifest")
	rootCmd.PersistentFlags().Bool("dry", false, "Print out the actions that would be taken without actually taking them")
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

func confirm(cmd *cobra.Command, msg string) bool {
	report.Println()
	report.PrintInfo(msg)

	reader := bufio.NewReader(os.Stdin)
	for {
		s, _ := reader.ReadString('\n')
		s = strings.ToLower(strings.TrimSuffix(s, "\n"))

		if strings.Compare(s, "n") == 0 {
			return false
		} else if strings.Compare(s, "y") == 0 {
			break
		} else {
			report.PrintInfo(msg)
		}
	}

	return true
}
