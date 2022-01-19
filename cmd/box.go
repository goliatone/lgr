package cmd

import (
	"fmt"

	"github.com/goliatone/lgr/pkg/widgets"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type boxOptions struct {
	Style     string
	Alignment string
	Hpad      int
	Vpad      int
}

var bo *boxOptions

func init() {
	bo = &boxOptions{}

	// boxCmd.Flags().IntVarP(&o.Total, "total", "t", 100, "Total to calculate progress.")
	boxCmd.Flags().StringVar(&bo.Style, "style", "double", "Border style: single, double, round, x, classic")
	boxCmd.Flags().StringVarP(&bo.Alignment, "alignment", "a", "center", "Box screen alignment: right, center, left.")
	boxCmd.Flags().IntVar(&bo.Hpad, "h-pad", 3, "Horizontal padding.")
	boxCmd.Flags().IntVar(&bo.Vpad, "v-pad", 2, "Vertical padding.")

	rootCmd.AddCommand(boxCmd)
}

var boxCmd = &cobra.Command{
	Use:   "box",
	Short: "prints a box",
	Long:  `print a box with the given context.`,
	Example: `

	`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		width, _, err := term.GetSize(0)
		if err != nil {
			return err
		}

		b := widgets.Box{
			Content:   args[0],
			Style:     bo.Style,
			Hpad:      bo.Hpad,
			Vpad:      bo.Vpad,
			ScreenW:   width,
			Alignment: bo.Alignment,
		}

		fmt.Println(b)

		return nil
	},
}
