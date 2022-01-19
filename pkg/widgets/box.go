package widgets

import (
	"strings"

	"github.com/jwalton/gchalk"
)

const defaultStyle = "single"

//BoxStyle holds box style
type BoxStyle struct {
	TopLeft     string
	Horizontal  string
	TopRight    string
	Vertical    string
	BottomRight string
	BottomLeft  string
}

func BoxStyleFromChar(char string) BoxStyle {
	return BoxStyle{
		TopRight:    char,
		TopLeft:     char,
		BottomRight: char,
		BottomLeft:  char,
		Horizontal:  char,
		Vertical:    char,
	}
}

//Box holds box info
type Box struct {
	Width     int
	MaxWidth  int
	Height    int
	Title     string
	Content   string
	Style     string
	Hpad      int
	Vpad      int
	ScreenW   int
	ScreenH   int
	Alignment string
}

var (
	boxStyles map[string]BoxStyle = map[string]BoxStyle{
		"single": {
			TopRight:    "┐",
			TopLeft:     "┌",
			BottomRight: "┘",
			BottomLeft:  "└",
			Horizontal:  "─",
			Vertical:    "│",
		},
		"double": {
			TopRight:    "╗",
			TopLeft:     "╔",
			BottomRight: "╝",
			BottomLeft:  "╚",
			Horizontal:  "═",
			Vertical:    "║",
		},
		"round": {
			TopRight:    "╮",
			TopLeft:     "╭",
			BottomRight: "╯",
			BottomLeft:  "╰",
			Horizontal:  "─",
			Vertical:    "│",
		},
		"x": {
			TopRight:    "+",
			TopLeft:     "+",
			BottomRight: "+",
			BottomLeft:  "+",
			Horizontal:  " ",
			Vertical:    " ",
		},
		"classic": {
			TopRight:    "+",
			TopLeft:     "+",
			BottomRight: "+",
			BottomLeft:  "+",
			Horizontal:  "-",
			Vertical:    "|",
		},
	}

	boxStyleAlias map[string]string = map[string]string{
		"rounded": "round",
	}
)

func (b Box) String() string {
	style := b.getStyle()

	lines := strings.Split(b.Content, "\\n")
	linesInfo := longestLine(lines)

	l := linesInfo.longest + 2 + (b.Hpad * 2)

	bar := strings.Repeat(style.Horizontal, l)
	tbar := style.TopLeft + bar + style.TopRight
	bbar := style.BottomLeft + bar + style.BottomRight

	padding := b.getVerticalPadding(l)
	hpadding := b.getHorizontalPadding(l)

	text := b.getText(l, linesInfo)

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

func (b Box) getStyle() BoxStyle {
	styleName := b.Style

	if _, ok := boxStyleAlias[styleName]; ok {
		styleName = boxStyleAlias[styleName]
	}

	if s, ok := boxStyles[styleName]; ok {
		return s
	}

	if len(styleName) == 1 {
		return BoxStyleFromChar(styleName)
	}

	return boxStyles[defaultStyle]
}

func (b Box) getHorizontalPadding(len int) string {
	switch b.Alignment {
	case "right":
		return strings.Repeat(" ", (b.ScreenW - len - 2))
	case "center":
		return strings.Repeat(" ", (b.ScreenW-len)/2)
	case "left":
		return ""
	default:
		return strings.Repeat(" ", (b.ScreenW-len)/2)
	}
}

func (b Box) getVerticalPadding(len int) []string {
	pad := strings.Repeat(" ", len)
	sep := b.getStyle().Vertical
	var lines = make([]string, 0, b.Vpad)
	for i := 0; i < b.Vpad; i++ {
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

func (b Box) getText(length int, lines textLines) []string {
	var s []string
	l := lines.lines
	sep := b.getStyle().Vertical

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
