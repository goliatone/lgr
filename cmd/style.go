package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(styleCmd)
}

var styleCmd = &cobra.Command{
	Use:     "style [text]",
	Aliases: []string{"s"},
	Short:   "Style given text",
	Example: `
	lgr -m bg-red -m grey -m bold "and italic and magenta"
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		opts.Heading = ""
		opts.NoTimestamp = true
		handleInput("", args)
	},
}
