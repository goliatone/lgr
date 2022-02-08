package box

import (
	"fmt"
	"strings"

	"github.com/jwalton/gchalk"
)

const defaultTemplate = "single"

//Template holds box style
type Template struct {
	TopLeft     string
	Horizontal  string
	TopRight    string
	Vertical    string
	BottomRight string
	BottomLeft  string
}

//TemplateFromChar generates a template from a single character
func TemplateFromChar(char string) Template {
	return Template{
		TopRight:    char,
		TopLeft:     char,
		BottomRight: char,
		BottomLeft:  char,
		Horizontal:  char,
		Vertical:    char,
	}
}

//Widget holds box info
type Widget struct {
	Width     int
	MaxWidth  int
	Height    int
	Title     string
	Content   string
	Template  string
	Hpad      int
	Vpad      int
	ScreenW   int
	ScreenH   int
	Alignment string
}

//New creates a new Widget instance
func New(options ...Option) *Widget {
	w := &Widget{}
	for _, option := range options {
		option(w)
	}
	return w
}

//SetContent updates content
func (w *Widget) SetContent(content string) *Widget {
	w.Content = content
	return w
}

//Render will render the widget
func (w *Widget) Render() {
	fmt.Println(w.String())
}

func (w *Widget) String() string {
	style := w.getTemplate()

	lines := strings.Split(w.Content, "\\n")
	linesInfo := longestLine(lines)

	l := linesInfo.longest + 2 + (w.Hpad * 2)

	bar := strings.Repeat(style.Horizontal, l)
	tbar := style.TopLeft + bar + style.TopRight
	bbar := style.BottomLeft + bar + style.BottomRight

	padding := w.getVerticalPadding(l)
	hpadding := w.getHorizontalPadding(l)

	text := w.getText(l, linesInfo)

	var sb strings.Builder

	sb.WriteString(hpadding + tbar)
	sb.WriteString("\n")

	for _, e := range padding {
		sb.WriteString(hpadding + e + "\n")
	}

	for i, e := range text {
		if i == 0 {
			s := gchalk.WithBrightWhite().WithBold()
			sb.WriteString(hpadding + s.Paint(e) + "\n")
		} else {
			sb.WriteString(hpadding + e + "\n")
		}
	}

	for _, e := range padding {
		sb.WriteString(hpadding + e + "\n")
	}

	sb.WriteString(hpadding + bbar)
	sb.WriteString("\n")

	return sb.String()
}

func (w *Widget) getTemplate() Template {
	styleName := w.Template

	if _, ok := boxTemplateAlias[styleName]; ok {
		styleName = boxTemplateAlias[styleName]
	}

	if s, ok := boxTemplates[styleName]; ok {
		return s
	}

	if len(styleName) == 1 {
		return TemplateFromChar(styleName)
	}

	return boxTemplates[defaultTemplate]
}

func (w *Widget) getHorizontalPadding(len int) string {
	switch w.Alignment {
	case "right":
		return strings.Repeat(" ", (w.ScreenW - len - 2))
	case "center":
		return strings.Repeat(" ", (w.ScreenW-len)/2)
	case "left":
		return ""
	default:
		return strings.Repeat(" ", (w.ScreenW-len)/2)
	}
}

func (w *Widget) getVerticalPadding(len int) []string {
	pad := strings.Repeat(" ", len)
	sep := w.getTemplate().Vertical
	var lines = make([]string, 0, w.Vpad)
	for i := 0; i < w.Vpad; i++ {
		lines = append(lines, (sep + pad + sep))
	}
	return lines
}

func longestLine(lines []string) textLines {
	t := textLines{
		longest: 0,
	}

	for _, line := range lines {
		l := 0
		for _, c := range line {
			l += runWidth(c)
		}
		t.lines = append(t.lines, lineInfo{
			len:  l,
			text: line,
		})

		if l > t.longest {
			//we always want the longest line to be even so that
			//padding works
			if l%2 != 0 {
				l++
			}
			t.longest = l
		}
	}
	return t
}

func runWidth(r rune) int {
	return 1
}

func (w *Widget) getText(length int, lines textLines) []string {
	var s []string
	l := lines.lines
	sep := w.getTemplate().Vertical

	for _, line := range l {

		plen := (length - line.len) / 2
		rpad := strings.Repeat(" ", plen)
		lpad := rpad

		if line.len%2 != 0 {
			lpad += " "
		}

		text := sep + rpad + line.text + lpad + sep

		s = append(s, text)
	}
	return s
}

type lineInfo struct {
	len  int
	text string
}

type textLines struct {
	longest int
	lines   []lineInfo
}
