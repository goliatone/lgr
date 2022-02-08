package bar

import "github.com/redmask-hb/GoSimplePrint/goPrint"

//Widget hold progress bar options
type Widget struct {
	Total           int
	Update          int
	Title           string
	DoneNotice      string
	GraphChar       string
	BackgroundChar  string
	HidePercent     bool
	HideRatio       bool
	HideProgressBar bool
	Colors          *Colors
}

//Colors has color options for the bar
type Colors struct {
	Graph   int
	Back    int
	Ratio   int
	Percent int
	Notice  int
}

//New will render the progress bar
//TODO: move to widgets
func New(options ...Option) *Widget {

	w := &Widget{}

	for _, option := range options {
		option(w)
	}
	return w
}

//NewWithDefaults creates a widget with default options
func NewWithDefaults() *Widget {
	return New(WithDefaultOptions())
}

//SetUpdate will set the Update value
func (w *Widget) SetUpdate(i int) *Widget {
	w.Update = i
	return w
}

//Render renders the bar
func (w *Widget) Render() {
	bar := goPrint.NewBar(w.Total)

	if w.Title != "" {
		bar.SetNotice(w.Title)
	}

	if w.HidePercent {
		bar.HidePercent()
	}

	if w.HideProgressBar {
		bar.HideProgressBar()
	}

	if w.HideRatio || w.Total == 100 {
		bar.HideRatio()
	}

	if w.Colors != nil {
		bar.SetColor(goPrint.BarColor{
			Graph:   w.Colors.Graph,
			Back:    w.Colors.Back,
			Ratio:   w.Colors.Ratio,
			Percent: w.Colors.Percent,
			Notice:  w.Colors.Notice,
		})
	}

	bar.SetGraph(w.GraphChar)
	bar.SetBackGraph(w.BackgroundChar)

	u := w.Update
	if u >= w.Total {
		u = w.Total
	}

	bar.PrintBar(u)

	if w.Update >= w.Total {
		bar.PrintEnd(w.DoneNotice)
	}
}
