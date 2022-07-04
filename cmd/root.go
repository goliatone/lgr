package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/goliatone/lgr/pkg/logging"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleLogStream(args)
	},
}

var errorExitCode = 1
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
	rootCmd.PersistentFlags().BoolVar(&opts.NoTimestamp, "no-timestamp", false, "do now show timestamp")

	opts.Modifiers = rootCmd.PersistentFlags().StringSliceP("modifier", "m", []string{}, "list of style modifiers")
}

//Execute exposes the root command execute method
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errorExitCode)
	}
}

func handleLogStream(args []string) error {

	parser := logging.JSONLineParser{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line, err := parser.Parse(scanner.Bytes())
		if err != nil {
			return err
		}
		//TODO: move to render
		message := fmt.Sprintf("%s   %s", line.Message, line.Fields)
		opts.Level = line.Level
		render.Print(message, opts)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func handleInput(level string, args []string) {
	opts.Level = level

	var scanner *bufio.Scanner

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		scanner = bufio.NewScanner(strings.NewReader(getBody(args)))
	}

	f := true
	for scanner.Scan() {
		body := scanner.Text()
		if f {
			body = indentOutput(body, opts.ShortHeading)
		}
		f = false
		render.Print(body, opts)
	}

	if opts.Level == "fatal" {
		os.Exit(errorExitCode)
	}
}

func getBody(args []string) string {
	if len(args) == 0 {
		return ""
	}
	//Get all strings as a single item
	return strings.Join(args, " ")
}
