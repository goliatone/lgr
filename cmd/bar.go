package cmd

import (
	"strconv"

	"github.com/goliatone/lgr/pkg/progress"
	"github.com/spf13/cobra"
)

var po *progress.Options

func init() {
	po = progress.DefaultOptions()

	barCmd.Flags().IntVarP(&po.Total, "total", "t", 100, "Total to calculate progress.")
	barCmd.Flags().StringVarP(&po.Title, "title", "T", "", "Title shown next to bar.")
	barCmd.Flags().StringVarP(&po.DoneNotice, "done", "d", "", "Message shown on completion.")
	barCmd.Flags().StringVarP(&po.GraphChar, "graph", "g", progress.DefaultGraphChar, "Character used to draw bar.")
	barCmd.Flags().StringVarP(&po.BackgroundChar, "back", "B", progress.DefaultBackgroundChar, "Character used to draw bar background.")
	barCmd.Flags().BoolVarP(&po.HidePercent, "percent", "p", false, "Hides the progress percent.")
	barCmd.Flags().BoolVarP(&po.HideRatio, "ratio", "r", true, "Hides the progress ratio.")
	barCmd.Flags().BoolVar(&po.HideProgressBar, "bar", false, "Hides the progress bar.")

	rootCmd.AddCommand(barCmd)
}

var barCmd = &cobra.Command{
	Use:   "bar",
	Short: "Prints a progress bar",
	Long:  `Print a progress bar and update on subsequent calls.`,
	Example: `
  lgr bar 10
  [█████#############################################] 10%

  lgr bar 100
  [██████████████████████████████████████████████████] 100%

  lgr bar -B _ 20 -g »
  [»»»»»»»»»»________________________________________] 20%
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		i, _ := strconv.Atoi(args[0])
		//TODO: handle error...

		po.Update = i

		progress.Bar(po)
	},
}
