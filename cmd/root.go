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
	Long: `lgr helps you colorize your script output

Styling:
To affect the printed output styling use
modifiers flags (-m,--modifier). You can use
multiple modifiers in a single invocation.

Style modifiers:
  * bold            * overline
  * dim             * reset
  * hidden          * strike through
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
  * --no-color: Set the global flag

Environment variable options:
  * LGR_COLOR: 			 --color
  * LGR_LEVEL: 			 --level
  * LGR_HEADING: 		 --heading
  * LGR_BOLD:   		 --bold
  * LGR_NO_COLOR:   	 --no-color
  * NO_COLOR:  			 --no-color
  * LGR_NO_NEWLINE: 	 --no-newline
  * LGR_SHORT_HEADLINES: --short-headlines
  * LGR_NO_TIMESTAMP: 	 --no-timestamp
  * LGR_TIME_FORMAT: 	 --time-format
  `,
	Example: `
	lgr error 'This is an error message'
	lgr -m bold -m underline -m bg-white -m magenta 'Hi there'
	echo "command output..." | lgr info
	`,
	Args: cobra.MinimumNArgs(0),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if _, ok := shortHeadings[cmd.CalledAs()]; ok {
			opts.ShortHeading = true
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 && args[0] != "" {
			handleInput(opts.Level, args)
			return nil
		}
		return handleLogStream(args)
	},
}

var errorExitCode = 1
var opts *render.Options

const maxBufferSize = 32 * 1024

func init() {
	opts = &render.Options{
		HeadingSuffix: " ",
	}

	pf := rootCmd.PersistentFlags()

	pf.StringVarP(
		&opts.Color,
		"color",
		"c",
		getEnv("LGR_COLOR", "neutral"),
		"line color",
	)

	pf.StringVarP(
		&opts.Level,
		"level",
		"l",
		getEnv("LGR_LEVEL", "trace"),
		"log level",
	)

	pf.StringVarP(
		&opts.Heading,
		"heading",
		"H",
		getEnv("LGR_HEADING", ""),
		"heading text",
	)

	pf.BoolVarP(
		&opts.Bold,
		"bold",
		"b",
		getEnvBool("LGR_BOLD", false),
		"bold style",
	)

	pf.BoolVar(
		&opts.NoColor,
		"no-color",
		getEnvBool("LGR_NO_COLOR", false),
		"disable color output",
	)
	pf.BoolVarP(
		&opts.NoNewline,
		"no-newline",
		"n",
		getEnvBool("LGR_NO_NEWLINE", false),
		"output not ended in newline",
	)

	pf.BoolVarP(
		&opts.ShortHeading,
		"short-headlines",
		"S",
		getEnvBool("LGR_SHORT_HEADLINES", false),
		"use short headings",
	)

	pf.BoolVar(
		&opts.NoTimestamp,
		"no-timestamp",
		getEnvBool("LGR_NO_TIMESTAMP", false),
		"do now show timestamp",
	)

	pf.StringVar(
		&opts.TimestampFormat,
		"time-format",
		getEnv("LGR_TIME_FORMAT", render.TimestampFormat),
		"timestamp format",
	)

	opts.Modifiers = pf.StringSliceP(
		"modifier",
		"m",
		[]string{},
		"list of style modifiers",
	)
}

//GetRoot returns the root command
func GetRoot() *cobra.Command {
	return rootCmd
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
	scanner.Buffer(make([]byte, maxBufferSize), maxBufferSize) // 32k

	i := 0
	for scanner.Scan() {
		b := scanner.Bytes()
		line, err := parser.Parse(b)
		if err != nil {
			//TODO: we need to implement a proper passthrough
			fmt.Println(string(b))
			continue
		}

		i++

		line.Line = i
		opts.Level = line.Level

		render.Print(line, opts)
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
		//TODO: maybe we also handle file paths? in which case we want to close handle
		scanner = bufio.NewScanner(strings.NewReader(getBody(args)))
	}

	i := 0
	for scanner.Scan() {
		body := scanner.Text()
		i++

		m := &logging.Message{
			Line:    i,
			Level:   level,
			Message: body,
		}

		render.Print(m, opts)
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

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getEnvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v == "true"
}
