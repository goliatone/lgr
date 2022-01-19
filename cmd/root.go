package cmd

import (
	"fmt"
	"os"

	"github.com/goliatone/lgr/pkg/render"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lgr [text|stdin]",
	Short: "Pretty print output to console",
	Long:  "lgr helps you colorize your script output",
	Example: `
	lgr error 'This is an error message'
	lgr -m bold -m underline -m bg-white -m magenta 'Hi there'
	echo "command output..." | lgr info

Styling:
  To affect the printed output styling use
  modifiers flags (-m,--modifier). You can use
  multiple modifiers in a single invocation.

  Style modifiers:
	* bold            * overline
	* dim             * reset
	* hidden          * strikethrough
	* inverse         * underline
	* italic

  Color modifiers:
	* black           * magenta
	* blue            * white
	* cyan            * yellow
	* green           * gray
	* red             * grey

  To set background color prepend "bg-" to any color
  e.g. bg-red for a red background.

  To use the bright version of a color prepend "hi-"
  to any color e.g. hi-red for bright red.

  Disable color:
    * NO_COLOR:   Set the environment variable
	* --no-color: Set the global flag
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

	rootCmd.PersistentFlags().StringVarP(&opts.Color, "color", "c", "neutral", "line color")
	rootCmd.PersistentFlags().StringVarP(&opts.Level, "level", "l", "trace", "log level")
	rootCmd.PersistentFlags().StringVarP(&opts.Heading, "heading", "H", "", "heading text")
	rootCmd.PersistentFlags().BoolVarP(&opts.Bold, "bold", "b", false, "bold style")
	rootCmd.PersistentFlags().BoolVar(&opts.NoColor, "no-color", false, "disable color output")
	rootCmd.PersistentFlags().BoolVarP(&opts.NoNewline, "no-newline", "n", false, "output not ended in newline")
	rootCmd.PersistentFlags().BoolVarP(&opts.ShortHeading, "short-headlines", "S", false, "use short headings")

	opts.Modifiers = rootCmd.PersistentFlags().StringSliceP("modifier", "m", []string{}, "list of style modifiers")
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
