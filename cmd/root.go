package cmd

import (
	"fmt"
	"os"

	"github.com/goliatone/lgr/pkg/render"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lgr",
	Short: "pretty print output to console",
	Long:  "lgr helps you colorize your script output",
	Example: `
	lgr error 'This is an error message'
	`,
	Args: cobra.MinimumNArgs(0),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if _, ok := shortHeadings[cmd.CalledAs()]; ok {
			opts.ShortHeading = true
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		handleInput("trace", args)
	},
}

var opts *render.Options

func init() {
	opts = &render.Options{}

	rootCmd.PersistentFlags().StringVarP(&opts.Color, "color", "c", "neutral", "Line color")
	rootCmd.PersistentFlags().StringVarP(&opts.Level, "level", "l", "trace", "Log level")
	rootCmd.PersistentFlags().BoolVarP(&opts.Bold, "bold", "b", false, "Bold style")
	rootCmd.PersistentFlags().BoolVarP(&opts.NoNewline, "no-newline", "n", false, "New line")
	rootCmd.PersistentFlags().BoolVarP(&opts.ShortHeading, "short-headlines", "S", false, "Short headings")

	opts.Modifiers = rootCmd.PersistentFlags().StringSliceP("modifier", "m", []string{}, "Style modifier")

}

//Execute exposes the root command execute method
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func handleInput(level string, args []string) {
	opts.Level = level
	//TODO: check if no body and no stdin then show usage?
	body := getBody(args)
	render.Print(body, opts)
}

func getBody(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}
