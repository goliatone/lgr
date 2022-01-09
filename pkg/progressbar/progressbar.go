package progressbar

import "github.com/redmask-hb/GoSimplePrint/goPrint"

//Options hold progress bar options
type Options struct {
	Total           int
	Update          int
	Title           string
	DoneNotice      string
	GraphChar       string
	BackgroundChar  string
	HidePercent     bool
	HideRatio       bool
	HideProgressBar bool
}

//Render will render the progress bar
func Render(o *Options) error {

	bar := goPrint.NewBar(o.Total)

	if o.Title != "" {
		bar.SetNotice(o.Title)
	}

	if o.HidePercent {
		bar.HidePercent()
	}

	if o.HideRatio || o.Total == 100 {
		bar.HideRatio()
	}

	if o.HideProgressBar {
		bar.HideProgressBar()
	}

	bar.SetColor(goPrint.BarColor{
		Graph:   32,
		Back:    32,
		Ratio:   32,
		Percent: 32,
		Notice:  32,
	})

	bar.SetGraph(o.GraphChar)
	bar.SetBackGraph(o.BackgroundChar)

	bar.PrintBar(o.Update)

	if o.Update == o.Total {
		if o.DoneNotice != "" {
			bar.PrintEnd(o.DoneNotice)
		}
	}

	return nil
}

//DefaultOptions returns an Options object with default values
func DefaultOptions() *Options {
	return &Options{
		Total:          100,
		Update:         0,
		DoneNotice:     "Task complete...",
		GraphChar:      "█",
		BackgroundChar: "#",
		// HideRatio:      true,
	}
}
