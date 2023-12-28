package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	defaultShell = "zsh"
)

var completionCmd = NewCompletionCmd(os.Stdout)

func init() {
	rootCmd.AddCommand(completionCmd)
}

func NewCompletionCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate completion files for the concord cli",
		RunE:  completionRun,
	}

	cmd.SetOut(out)

	cmd.Flags().String("shell", defaultShell, "desired shell to generate completions for")

	return cmd
}

func completionRun(cmd *cobra.Command, args []string) error {
	shell := cmd.Flags().Lookup("shell").Value.String()

	var err error
	switch {
	case strings.EqualFold(shell, "bash"):
		err = rootCmd.GenBashCompletion(os.Stdout)

	case strings.EqualFold(shell, "ps") ||
		strings.EqualFold(shell, "powershell") ||
		strings.EqualFold(shell, "power_shell"):
		err = rootCmd.GenPowerShellCompletion(os.Stdout)

	case strings.EqualFold(shell, "zsh"):
		err = rootCmd.GenZshCompletion(os.Stdout)

	default:
		err = fmt.Errorf("unsupported shell type: %s", shell)
	}

	if err != nil {
		return handleError(cmd, fmt.Errorf("completion: %w", err))
	}

	return nil
}
