package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(successCmd)
	rootCmd.AddCommand(failureCmd)
}

var successCmd = &cobra.Command{
	Use:     "success [text|stdin]",
	Aliases: []string{"ok", "OK"},
	Short:   "Show a success notice",
	// Long:    generateLong("DEBUG", "D"),
	Example: `
  lgr success 'This is a successful message'
  lgr ok 'This is a successful message'
  curl -v example.com | lgr ok
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleInput("success", args)
	},
}

var failureCmd = &cobra.Command{
	Use:     "failure [text|stdin]",
	Aliases: []string{"fail", "ko", "KO"},
	Short:   "Show a failure notice",
	// Long:    generateLong("INFO", "I"),
	Example: `
	lgr failure 'This is a failed message'
  lgr ko 'This is a failed message'
  curl -v example.com || lgr ko
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleInput("failure", args)
	},
}
