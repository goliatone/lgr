package cmd

import (
	"strconv"

	"github.com/goliatone/lgr/pkg/widgets/bar"
	"github.com/spf13/cobra"
)

var widget *bar.Widget

func init() {
	widget = bar.NewWithDefaults()

	barCmd.Flags().IntVarP(&widget.Total, "total", "t", 100, "Total to calculate progress.")
	barCmd.Flags().StringVarP(&widget.Title, "title", "T", "", "Title shown next to bar.")
	barCmd.Flags().StringVarP(&widget.DoneNotice, "done", "d", "", "Message shown on completion.")
	barCmd.Flags().StringVarP(&widget.GraphChar, "graph", "g", bar.DefaultGraphChar, "Character used to draw bar.")
	barCmd.Flags().StringVarP(&widget.BackgroundChar, "back", "B", bar.DefaultBackgroundChar, "Character used to draw bar background.")
	barCmd.Flags().BoolVarP(&widget.HidePercent, "percent", "p", false, "Hides the progress percent.")
	barCmd.Flags().BoolVarP(&widget.HideRatio, "ratio", "r", true, "Hides the progress ratio.")
	barCmd.Flags().BoolVar(&widget.HideProgressBar, "bar", false, "Hides the progress bar.")

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
	RunE: func(cmd *cobra.Command, args []string) error {

		i, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		widget.SetUpdate(i).Render()
		return nil
	},
}
