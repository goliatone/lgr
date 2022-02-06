package render

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jwalton/gchalk"
)

//Options holds print modifiers
type Options struct {
	Bold          bool
	Level         string
	Color         string
	NoColor       bool
	Heading       string
	ShortHeading  bool
	HeadingPrefix string
	NoNewline     bool
	Modifiers     *[]string
}

//Stylize will add stile to your body
func Stylize(body string, opts *Options) (string, string) {
	//Add heading
	heading := headings[opts.Level]

	if opts.ShortHeading {
		heading = headingShort[opts.Level]
	}

	if opts.Heading != "" {
		heading = opts.Heading
	}

	if heading != "" {
		if style, ok := headingStyle[opts.Level]; ok {
			heading = style.Paint(heading)
			heading += " "
		}

		if opts.HeadingPrefix != "" {
			heading = opts.HeadingPrefix + heading
		}
	}

	content := body

	if ok, _ := checkInput(os.Stdin); ok {
		if body != "" {
			content += " "
		}
		content = streamToString(os.Stdin)
	}

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

	//Check before we add asci chars and we cant check
	if opts.NoNewline && !strings.HasSuffix(content, " ") {
		content = content + " "
	}

	content = style.Paint(content)

	if !opts.NoNewline {
		content = content + "\n"
	}

	return heading, content
}

//Print will render content to stdout
func Print(body string, opts *Options) {
	heading, content := Stylize(body, opts)
	//TODO: use writer
	// fmt.Fprint(w io.Writer, a ...interface{})

	fmt.Printf("%s%s", heading, content)
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
