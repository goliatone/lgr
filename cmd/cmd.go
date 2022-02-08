package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/goliatone/lgr/pkg/widgets"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func init() {
	rootCmd.AddCommand(execCmd)

}

var execCmd = &cobra.Command{
	Use:     "exec [text|stdin]",
	Aliases: []string{"ok", "OK"},
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

		opts.ShortHeading = false

		label := "executing: " + strings.Join(args, " ")

		name, params := makeExecParams(args)
		run := exec.Command(name, params...)
		run.Env = os.Environ()

		var out bytes.Buffer
		run.Stdout = &out
		run.Stderr = &out
		s := widgets.NewSpinner(label)
		s.MaxWidth = getMaxScreenWidth(3)
		s.Frames = widgets.FramesBarHorizontal
		s.Output = os.Stdout

		s.Start()

		err := run.Run()

		if err != nil {
			s.Stop()
			handleInput("failure", []string{"error " + label})

			opts.WithIndent()
			content := indentOutput(err.Error()+"\n"+out.String(), opts.ShortHeading)
			handleInput("error", []string{content})

			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			}
			return
		}

		s.Stop()

		handleInput("success", []string{"success " + label})

		content := indentOutput(out.String(), opts.ShortHeading)

		if content == "" {
			return
		}

		opts.WithIndent()
		handleInput("info", []string{content})
	},
}

func makeExecParams(args []string) (string, []string) {

	return args[0], args[1:]
}

func indentOutput(input string, short bool) string {
	indentLength := 7
	if short {
		indentLength = 3
	}
	indent := fmt.Sprintf("%*s", indentLength, "")

	out := []string{}

	lines := strings.Split(strings.TrimRight(input, "\n"), "\n")

	lineDigits := countDigits(len(lines))
	for i, l := range lines {
		l = strings.TrimLeft(l, " ")
		if i == 0 {
			out = append(out, fmt.Sprintf("%s%s", fmt.Sprintf("%*s", lineDigits-1, ""), l))
		} else {
			// out = append(out, fmt.Sprintf("   [%0*d] %s", lineDigits, i, l))
			out = append(out, fmt.Sprintf("%s[%0*d] %s", indent, lineDigits, i, l))
		}
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
