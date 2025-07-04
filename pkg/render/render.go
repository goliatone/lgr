package render

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/goliatone/lgr/pkg/logging"
	"github.com/jwalton/gchalk"
)

// Options holds print modifiers
type Options struct {
	Bold            bool
	Level           string
	Color           string
	NoColor         bool
	Heading         string
	ShortHeading    bool
	HeadingPrefix   string
	HeadingSuffix   string
	NoNewline       bool
	NoTimestamp     bool
	Modifiers       *[]string
	TimestampFormat string
	MaxBufferSize   int
	Writer          io.Writer
	Filters         *[]string
}

const clear = "\x1b[0m"

// IndentationChar is the character used for indentation
var IndentationChar string = " └─"

// TimestampFormat is the default timestamp format
var TimestampFormat = "01-02-2006 15:04:05.000000"

// WithIndent sets heading with indent option
func (o *Options) WithIndent() *Options {
	o.HeadingPrefix = IndentationChar
	return o
}

// WithHeadingSuffix sets heading suffix
func (o *Options) WithHeadingSuffix(s string) *Options {
	o.HeadingSuffix = s
	return o
}

// HasIndent returns true if heading has indent
func (o *Options) HasIndent() bool {
	return o.HeadingPrefix == IndentationChar
}

// Stylize will add stile to your body
// TODO: use Message interface instead of struct to prevent cyclic deps
func Stylize(msg *logging.Message, opts *Options) (string, string) {

	if msg.HasFields() {
		fieldTemplate := getFieldTemplate(opts)
		msg.WithFieldTemplate(fieldTemplate)
		clearSeq := getClearSequence(opts)
		msg.Message = fmt.Sprintf("%s%s%s%s", msg.Message, clearSeq, "\t", msg.Fields)
	}

	body := msg.Message

	//Add heading
	heading := getHeading(opts)
	heading = styleHeading(heading, opts)

	if msg.Caller != "" {
		heading = heading + "<" + msg.Caller + "> "
	}

	now := msg.GetTimestampOrNow()
	ts := now.Format(opts.TimestampFormat)
	ts = applyTimestampStyle(ts, opts)

	if opts.NoTimestamp != true {
		if opts.HasIndent() {
			ts = strings.Repeat(" ", utf8.RuneCountInString(opts.TimestampFormat))
		}
		heading = fmt.Sprintf("%s %s", ts, heading)
	}

	content := body

	if !opts.NoColor {
		style, err := gchalk.WithStyle(opts.Color)

		if opts.Modifiers != nil && len(*opts.Modifiers) > 0 {
			m := normalizeStyles(*opts.Modifiers...)
			style, err = style.WithStyle(m...)
		} else if mods, ok := modifiers[opts.Level]; ok {
			m := normalizeStyles(mods...)
			style, err = style.WithStyle(m...)
		}

		if err != nil {
			//TODO: check for tty theme
			style = gchalk.WithBrightWhite()
		}

		if opts.Bold {
			style = style.WithBold()
		}

		content = style.Paint(content)
	} else {
		// When colors are disabled, we might still want to handle bold differently
		// For now, we'll just use the plain content
	}

	if opts.NoNewline && !strings.HasSuffix(content, " ") {
		content = content + " "
	}

	if msg.Stacktrace != "" {
		content += "\n\t" + strings.ReplaceAll(msg.Stacktrace, "\n", "\n\t")
	}

	if !opts.NoNewline {
		content = content + "\n"
	}

	return heading, content
}

// Print will render content to stdout
func Print(msg *logging.Message, opts *Options) {
	heading, content := Stylize(msg, opts)
	fmt.Fprintf(opts.Writer, "%s%s", heading, content)
}

func checkInput(f *os.File) (bool, error) {
	fi, err := f.Stat()
	if err != nil {
		return false, fmt.Errorf("failed reading stdin: %w", err)
	}
	ok := fi.Size() > 0
	return ok, nil
}

func streamToString(s io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(s)
	return buf.String()
}

func normalizeStyles(styles ...string) []string {
	o := make([]string, len(styles))
	for i, m := range styles {
		v, ok := styleMap[m]
		if ok {
			o[i] = v
		} else {
			o[i] = m
		}
	}
	return o
}

func getFieldTemplate(opts *Options) string {
	if opts.NoColor {
		return "%s=%s"
	}
	return "\x1b[38;5;244m%s\x1b[0m=%s"
}

func getClearSequence(opts *Options) string {
	if opts.NoColor {
		return ""
	}
	return clear
}

func applyTimestampStyle(ts string, opts *Options) string {
	if opts.NoColor {
		return ts
	}

	if style, ok := elementStyle["timestamp"]; ok {
		return style.Paint(ts)
	}

	return ts
}

func getHeading(opts *Options) string {
	if opts.Heading != "" {
		return opts.Heading
	}
	heading := headings[opts.Level]

	if opts.ShortHeading {
		heading = headingShort[opts.Level]
	}
	return heading
}

func styleHeading(heading string, opts *Options) string {
	if heading == "" {
		return heading
	}

	if !opts.NoColor {
		if style, ok := headingStyle[opts.Level]; ok {
			heading = style.Paint(heading)
		}
	}

	heading += opts.HeadingSuffix

	if opts.HasIndent() {
		heading = opts.HeadingPrefix + heading
	}

	return heading
}
