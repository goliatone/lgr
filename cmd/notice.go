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
  lgr success 'This is a debug message'
  lgr D 'This is a debug message'
  curl -v example.com | lgr debug
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		handleInput("success", args)
	},
}

var failureCmd = &cobra.Command{
	Use:     "failure [text|stdin]",
	Aliases: []string{"fail", "ko", "KO"},
	Short:   "Show a failure notice",
	// Long:    generateLong("INFO", "I"),
	Example: `
	lgr info 'This is a debug message'
	lgr I 'This is a debug message'
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		handleInput("failure", args)
	},
}
