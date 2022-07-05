package render

import "github.com/jwalton/gchalk"

var elementStyle = map[string]*gchalk.Builder{
	"timestamp": gchalk.WithGray().WithBold(),
	"body":      gchalk.WithGray(),
}

var headings = map[string]string{
	"trace":   "",
	"debug":   "[DEBUG]",
	"info":    "[INFO ]",
	"warn":    "[WARN ]",
	"error":   "[ERROR]",
	"fatal":   "[FATAL]",
	"success": " 👍",
	"failure": " 👎",
}

var headingShort = map[string]string{
	"trace":   "",
	"debug":   "[D]",
	"info":    "[I]",
	"warn":    "[W]",
	"error":   "[E]",
	"fatal":   "[F]",
	"success": "[✔]",
	"failure": "[✖]",
}

//Default heading headingStyle
var headingStyle = map[string]*gchalk.Builder{
	"trace":   gchalk.WithBrightWhite(),
	"debug":   gchalk.WithBrightWhite().WithBold(),
	"info":    gchalk.WithBrightCyan().WithBold(),
	"warn":    gchalk.WithBrightYellow().WithBold(),
	"error":   gchalk.WithRed().WithBold(),
	"fatal":   gchalk.WithBrightRed().WithBold(),
	"success": gchalk.WithGreen().WithBold(),
	"failure": gchalk.WithRed().WithBold(),
}

//Default body style modifiers
var modifiers = map[string][]string{
	"fatal":   {"bg-red", "bold", "white"},
	"success": {"bold", "green"},
	"failure": {"bold", "red"},
}

var styleMap = map[string]string{
	//Style
	"bold":          "bold",
	"dim":           "dim",
	"hidden":        "hidden",
	"inverse":       "inverse",
	"italic":        "italic",
	"overline":      "overline",
	"reset":         "reset",
	"strikethrough": "strikethrough",
	"underline":     "underline",

	//Colors
	"black":   "black",
	"blue":    "blue",
	"cyan":    "cyan",
	"green":   "green",
	"magenta": "magenta",
	"red":     "red",
	"white":   "white",
	"yellow":  "yellow",
	"gray":    "gray",
	"grey":    "grey",

	"hi-black":   "brightBlack",
	"hi-blue":    "brightBlue",
	"hi-cyan":    "brightCyan",
	"hi-green":   "brightGreen",
	"hi-magenta": "brightMagenta",
	"hi-red":     "brightRed",
	"hi-white":   "brightWhite",
	"hi-yellow":  "brightYellow",

	"bg-black":   "bgBlack",
	"bg-blue":    "bgBlue",
	"bg-cyan":    "bgCyan",
	"bg-gray":    "bgGray",
	"bg-grey":    "bgGrey",
	"bg-green":   "bgGreen",
	"bg-magenta": "bgMagenta",
	"bg-red":     "bgRed",
	"bg-white":   "bgWhite",
	"bg-yellow":  "bgYellow",

	"bg-hi-black":   "bgBrightBlack",
	"bg-hi-blue":    "bgBrightBlue",
	"bg-hi-cyan":    "bgBrightCyan",
	"bg-hi-green":   "bgBrightGreen",
	"bg-hi-magenta": "bgBrightMagenta",
	"bg-hi-red":     "bgBrightRed",
	"bg-hi-white":   "bgBrightWhite",
	"bg-hi-yellow":  "bgBrightYellow",

	"hi-bg-black":   "bgBrightBlack",
	"hi-bg-blue":    "bgBrightBlue",
	"hi-bg-cyan":    "bgBrightCyan",
	"hi-bg-green":   "bgBrightGreen",
	"hi-bg-magenta": "bgBrightMagenta",
	"hi-bg-red":     "bgBrightRed",
	"hi-bg-white":   "bgBrightWhite",
	"hi-bg-yellow":  "bgBrightYellow",
}
