package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/goliatone/lgr/pkg/widgets"
	"github.com/spf13/cobra"
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

		opts.ShortHeading = true

		label := "executing " + strings.Join(args, " ")

		name, params := makeExecParams(args)
		run := exec.Command(name, params...)
		run.Env = os.Environ()

		var out bytes.Buffer
		run.Stdout = &out

		s := widgets.NewSpinner(label)

		s.Frames = widgets.FramesBarsVertical
		s.Output = os.Stdout

		s.Start()

		err := run.Run()

		if err != nil {
			s.Stop()
			handleInput("failure", []string{"error " + label})
			opts.HeadingPrefix = " └─"
			handleInput("error", []string{err.Error()})
			return
		}

		s.Stop()
		handleInput("success", []string{"success " + label})
		opts.HeadingPrefix = " └─"
		handleInput("info", []string{out.String()})
		// fmt.Printf("output: %q\n", out.String())
	},
}

func makeExecParams(args []string) (string, []string) {
	return args[0], args[1:]
}
