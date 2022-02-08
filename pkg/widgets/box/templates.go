package box

var (
	boxTemplates map[string]Template = map[string]Template{
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

	boxTemplateAlias map[string]string = map[string]string{
		"rounded": "round",
	}
)
