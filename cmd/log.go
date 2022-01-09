package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(fatalCmd)
	rootCmd.AddCommand(errorCmd)
	rootCmd.AddCommand(warnCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(debugCmd)
}

var shortHeadings = map[string]struct{}{
	"D": {},
	"I": {},
	"W": {},
	"E": {},
	"F": {},
}

var debugCmd = &cobra.Command{
	Use:     "debug",
	Aliases: []string{"D"},
	Short:   "Prepend [DEBUG] hading to output",
	Long:    "lgr is a terminal console",
	Example: `
	lgr debug 'This is a debug message'
	lgr D 'This is a debug message'
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		handleInput("debug", args)
	},
}

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"I"},
	Short:   "Prepend [INFO] hading to output",
	Long:    "lgr is a terminal console",
	Example: "lgr -l e 'This is my log message'",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		handleInput("info", args)
	},
}

var warnCmd = &cobra.Command{
	Use:     "warn",
	Aliases: []string{"W", "warning"},
	Short:   "Prepend [WARN] hading to output",
	Long:    "lgr is a terminal console",
	Example: "lgr -l e 'This is my log message'",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		handleInput("warn", args)
	},
}

var errorCmd = &cobra.Command{
	Use:     "error",
	Aliases: []string{"E", "err"},
	Short:   "Prepend [ERROR] hading to output",
	Long:    "lgr is a terminal console",
	Example: "lgr -l e 'This is my log message'",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		handleInput("error", args)
	},
}

var fatalCmd = &cobra.Command{
	Use:     "fatal",
	Aliases: []string{"F"},
	Short:   "Prepend [FATAL] hading to output",
	Long:    "lgr is a terminal console",
	Example: "lgr -l e 'This is my log message'",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		handleInput("fatal", args)
	},
}
