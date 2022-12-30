package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/goliatone/lgr/pkg/widgets/spinner"
	"github.com/jwalton/gchalk"
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
	RunE: cmdExec,
}

func cmdExec(cmd *cobra.Command, args []string) error {
	heading, args := getCmdArgs(os.Args[2:], args)

	label := heading + ": " + strings.Join(args, " ")

	// opts.ShortHeading = true

	var output strings.Builder
	output.WriteString(strings.Join(args, " ") + "\n")

	name, params := makeExecParams(args)

	run := exec.Command(name, params...)
	run.Env = os.Environ()
	run.Stderr = &output

	stdout, err := run.StdoutPipe()
	if err != nil {
		return err
	}

	//TODO: make configurable via flags
	widget := spinner.New(
		// spinner.WithLabel(label),
		spinner.WithMaxWidth(getMaxScreenWidth(3)),
		spinner.WithFrames(spinner.FramesSnake),
		spinner.WithOutput(os.Stdout),
		spinner.WithStyle(gchalk.WithBrightCyan()),
		spinner.WithFrameRate(time.Millisecond*60),
	)
	defer widget.Close()

	widget.SetLabel(label)

	done := stdoutScanner(heading, stdout, &output, widget)

	widget.Start()
	run.Start()

	<-done

	err = run.Wait()
	widget.Stop()

	if err != nil {
		handleInput("failure", []string{"error: " + err.Error()})

		opts.WithIndent()
		content := indentOutput(err.Error()+"\n"+output.String(), opts.ShortHeading)
		handleInput("error", []string{content})

		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		return nil
	}

	handleInput("success", []string{"success " + label})

	content := indentOutput(output.String(), opts.ShortHeading)
	if content == "" {
		return nil
	}

	opts.
		WithIndent().
		WithHeadingSuffix("")

	handleInput("info", []string{content})

	return nil
}

func stdoutScanner(heading string, stdout io.Reader, output *strings.Builder, widget *spinner.Widget) chan struct{} {
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	done := make(chan struct{})

	go func() {
		for scanner.Scan() {
			m := scanner.Text()
			output.WriteString(m + "\n")
			widget.UpdateLabel(heading + " " + m)
		}
		done <- struct{}{}
	}()

	return done
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
