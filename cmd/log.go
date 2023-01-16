package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {

	fatalCmd.Flags().IntVar(&errorExitCode, "exit-code", 1, "error exit code")

	rootCmd.AddCommand(fatalCmd)
	rootCmd.AddCommand(errorCmd)
	rootCmd.AddCommand(warnCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(debugCmd)
}

var shortHeadings = map[string]struct{}{
	"D":  {},
	"I":  {},
	"W":  {},
	"E":  {},
	"F":  {},
	"ok": {},
	"OK": {},
	"ko": {},
	"KO": {},
}

var debugCmd = &cobra.Command{
	Use:     "debug [text|stdin]",
	Aliases: []string{"D"},
	Short:   "Prepend [DEBUG] to text",
	Long:    generateLong("DEBUG", "D"),
	Example: `
  lgr debug 'This is a debug message'
  lgr D 'This is a debug message'
  curl -v example.com | lgr debug
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleInput("debug", args)
	},
}

var infoCmd = &cobra.Command{
	Use:     "info [text|stdin]",
	Aliases: []string{"I"},
	Short:   "Prepend [INFO] to text",
	Long:    generateLong("INFO", "I"),
	Example: `
	lgr info 'This is a debug message'
	lgr I 'This is a debug message'
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleInput("info", args)
	},
}

var warnCmd = &cobra.Command{
	Use:     "warn [text|stdin]",
	Aliases: []string{"W", "warning"},
	Short:   "Prepend [WARN] to text",
	Long:    generateLong("WARN", "W"),
	Example: `
	lgr warn 'This is a debug message'
	lgr W 'This is a debug message'
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleInput("warn", args)
	},
}

var errorCmd = &cobra.Command{
	Use:     "error [text|stdin]",
	Aliases: []string{"E", "err"},
	Short:   "Prepend [ERROR] to text",
	Long:    generateLong("ERROR", "E"),
	Example: `
	lgr error 'This is a debug message'
	lgr E 'This is a debug message'
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleInput("error", args)
	},
}

var fatalCmd = &cobra.Command{
	Use:     "fatal [text|stdin]",
	Aliases: []string{"F"},
	Short:   "Prepend [FATAL] to text",
	Long:    generateLong("FATAL", "F"),
	Example: `
	lgr fatal 'This is a debug message'
	lgr F 'This is a debug message'
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleInput("fatal", args)
	},
}

func generateLong(heading, short string) string {
	return fmt.Sprintf(`Prepend [%s] to text.
If the "%s" alias is used the heading used is [%s].

Input can be a string passed as the call argument.

If the command is called with text piped then it will
be appended to the output.`, heading, short, short)
}
