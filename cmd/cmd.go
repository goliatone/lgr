package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/goliatone/lgr/pkg/widgets/spinner"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func init() {
	rootCmd.AddCommand(execCmd)

}

const defaultExecHeading = "executing"

var execCmd = &cobra.Command{
	Use:     "exec [text|stdin]",
	Aliases: []string{"cmd"},
	Short:   "Execute command and show either success or failure notice",
	Long: `Execute command and show either success or failure notice.
If the process exits without error (exit code 0) then we
print a success message, any other exit code will print a
failure message.
	`,
	Example: `
  lgr exec 'This is a successful message'
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		heading, args := getCmdArgs(os.Args[2:], args)

		label := heading + ": " + strings.Join(args, " ")

		name, params := makeExecParams(args)
		run := exec.Command(name, params...)
		run.Env = os.Environ()

		var out bytes.Buffer
		run.Stdout = &out
		run.Stderr = &out

		//TODO: make configurable via flags
		widget := spinner.New(
			spinner.WithLabel(label),
			spinner.WithMaxWidth(getMaxScreenWidth(3)),
			spinner.WithFrames(spinner.FramesBarHorizontal),
			spinner.WithOutput(os.Stdout),
		)

		widget.Start()

		err := run.Run()

		if err != nil {
			widget.Stop()
			handleInput("failure", []string{"error " + label})

			opts.WithIndent()
			content := indentOutput(err.Error()+"\n"+out.String(), opts.ShortHeading)
			handleInput("error", []string{content})

			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			}
			return
		}

		widget.Stop()

		handleInput("success", []string{"success " + label})

		content := indentOutput(out.String(), opts.ShortHeading)

		if content == "" {
			return
		}

		opts.
			WithIndent().
			WithHeadingSuffix("")

		handleInput("info", []string{content})
	},
}

func getCmdArgs(o, c []string) (string, []string) {

	for i := len(o) - 1; i >= 0; i-- {
		arg := o[i]
		if arg == "--" {

			if i == 0 {
				return defaultExecHeading, o[i+1:]
			}
			r := o[i+1:]
			return strings.Join(c[0:len(r)], " "), r
		}
	}

	return defaultExecHeading, c
}

func makeExecParams(args []string) (string, []string) {
	return args[0], args[1:]
}

func indentOutput(input string, short bool) string {
	out := []string{}

	lines := strings.Split(strings.TrimRight(input, "\n"), "\n")

	lineDigits := countDigits(len(lines))
	for i, l := range lines {
		l = strings.TrimLeft(l, " ")
		out = append(out, fmt.Sprintf("[%0*d] %s", lineDigits, i+1, l))
	}

	return strings.Join(out, "\n")
}

func countDigits(i int) int {
	if i < 10 {
		return 1
	}
	return 1 + countDigits(i/10)
}

func getMaxScreenWidth(padding int) int {
	maxLineWidth := 80

	if term.IsTerminal(0) {
		maxLineWidth, _, _ = term.GetSize(0)
	}

	return maxLineWidth - padding
}
