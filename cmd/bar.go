package cmd

import (
	"strconv"

	"github.com/goliatone/lgr/pkg/progressbar"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(barCmd)
}

var barCmd = &cobra.Command{
	Use:     "bar",
	Short:   "print progress bar",
	Long:    "plog is a terminal console",
	Example: "plog bar 10",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		i, _ := strconv.Atoi(args[0])

		o := progressbar.DefaultOptions()
		o.Update = i
		// o.Title = "Building... "

		progressbar.Render(o)
	},
}
