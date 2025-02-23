package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gomicro/concord/client"
	"github.com/gomicro/concord/config"
	"github.com/gomicro/concord/report"
	"github.com/gomicro/scribe"
	"github.com/spf13/cobra"
)

var scrb scribe.Scriber

func init() {
	cobra.OnInitialize(initEnvs)

	rootCmd.PersistentFlags().StringP("file", "f", "concord.yml", "Path to a file containing a manifest")
	rootCmd.PersistentFlags().Bool("dry", false, "Print out the actions that would be taken without actually taking them")
	rootCmd.PersistentFlags().Bool("force", false, "Force the action to be taken without prompting for confirmation")

	t := &scribe.Theme{
		Describe: func(s string) string {
			// Cyan
			return fmt.Sprintf("\033[36m%s\033[0m", s)
		},
		Done: scribe.NoopDecorator,
	}

	scrb = scribe.NewScribe(os.Stdout, t)
}

func initEnvs() {
}

var rootCmd = &cobra.Command{
	Use:   "concord",
	Short: "concord is a tool to manage your Github repositories",
}

func Execute() {
	c, err := config.ParseFromFile()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	tkn := os.Getenv("GITHUB_TOKEN")

	if tkn != "" {
		c.Github.Token = tkn
	}

	ctx, err := client.WithClient(context.Background(), c.Github.Token)
	if err != nil && !errors.Is(err, client.ErrTokenEmpty) {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	err = rootCmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func handleError(c *cobra.Command, err error) error {
	c.SilenceUsage = true
	return err
}

func confirm(cmd *cobra.Command, msg string) bool {
	if strings.EqualFold(cmd.Flags().Lookup("force").Value.String(), "true") {
		return true
	}

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
