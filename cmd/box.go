package cmd

import (
	"fmt"

	"github.com/goliatone/lgr/pkg/widgets/box"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type boxOptions struct {
	Style     string
	Alignment string
	Hpad      int
	Vpad      int
	ScreenW   int
}

var boxWidget *box.Widget

func init() {
	boxWidget = box.New()

	// boxCmd.Flags().IntVarP(&o.Total, "total", "t", 100, "Total to calculate progress.")
	boxCmd.Flags().StringVar(&boxWidget.Template, "style", "double", "Border style: single, double, round, x, classic")
	boxCmd.Flags().StringVarP(&boxWidget.Alignment, "alignment", "a", "center", "Box screen alignment: right, center, left.")
	boxCmd.Flags().IntVar(&boxWidget.Hpad, "h-pad", 3, "Horizontal padding.")
	boxCmd.Flags().IntVar(&boxWidget.Vpad, "v-pad", 2, "Vertical padding.")
	boxCmd.Flags().IntVar(&boxWidget.ScreenW, "columns", 80, "Screen colum width e.g $COLUMNS.")

	rootCmd.AddCommand(boxCmd)
}

var boxCmd = &cobra.Command{
	Use:   "box",
	Short: "prints a box",
	Long:  `print a box with the given context.`,
	Example: `

	`,
	// Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if boxWidget.ScreenW == 0 && term.IsTerminal(0) {
			boxWidget.ScreenW, _, err = term.GetSize(0)
			if err != nil {
				return fmt.Errorf("run command: %w", err)
			}
		}

		boxWidget.SetContent(getBody(args)).Render()

		return nil
	},
}
