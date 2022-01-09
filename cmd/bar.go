package cmd

import (
	"strconv"

	"github.com/goliatone/lgr/pkg/progressbar"
	"github.com/spf13/cobra"
)

var o *progressbar.Options

func init() {
	o = progressbar.DefaultOptions()

	barCmd.Flags().IntVarP(&o.Total, "total", "t", 100, "Total to calculate progress.")
	barCmd.Flags().StringVarP(&o.Title, "title", "T", "", "Title shown next to bar.")
	barCmd.Flags().StringVarP(&o.DoneNotice, "done", "d", "", "Message shown on completion.")
	barCmd.Flags().StringVarP(&o.GraphChar, "graph", "g", progressbar.DefaultGraphChar, "Character used to draw bar.")
	barCmd.Flags().StringVarP(&o.BackgroundChar, "back", "B", progressbar.DefaultBackgroundChar, "Character used to draw bar background.")
	barCmd.Flags().BoolVarP(&o.HidePercent, "percent", "p", false, "Hides the progress percent.")
	barCmd.Flags().BoolVarP(&o.HideRatio, "ratio", "r", true, "Hides the progress ratio.")
	barCmd.Flags().BoolVar(&o.HideProgressBar, "bar", false, "Hides the progress bar.")

	rootCmd.AddCommand(barCmd)
}

var barCmd = &cobra.Command{
	Use:   "bar",
	Short: "prints a progress bar",
	Long: `bar will print a progress bar and update
on subsequent calls.`,
	Example: "lgr bar 10",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		i, _ := strconv.Atoi(args[0])
		//TODO: handle error...

		o.Update = i

		progressbar.Render(o)
	},
}
