package bar

const (
	//DefaultDoneNotice is the default message shown at completion
	DefaultDoneNotice = "Task complete..."
	//DefaultGraphChar is default char for bar graphic
	DefaultGraphChar string = "█"
	//DefaultBackgroundChar is default char for bar background graphic
	DefaultBackgroundChar string = "#"
)

//Option is a widget setter
type Option func(*Widget)

//WithDefaultOptions sets default options
func WithDefaultOptions() Option {
	return func(w *Widget) {
		w.Update = 0
		w.HideRatio = true
		w.DoneNotice = DefaultDoneNotice
		w.Total = 100
		w.GraphChar = DefaultGraphChar
		w.BackgroundChar = DefaultBackgroundChar
	}
}

//WithTitle set title value
func WithTitle(value string) Option {
	return func(w *Widget) {
		w.Title = value
	}
}

//WithDoneNotice set title value
func WithDoneNotice(value string) Option {
	return func(w *Widget) {
		w.DoneNotice = value
	}
}

//WithHidePercent set title value
func WithHidePercent(value bool) Option {
	return func(w *Widget) {
		w.HidePercent = value

	}
}

// WithHideRatio set hide ratio
func WithHideRatio(value bool) Option {
	return func(w *Widget) {
		w.HideRatio = value
	}
}

// WithHideProgressBar set hide ratio
func WithHideProgressBar(value bool) Option {
	return func(w *Widget) {
		w.HideProgressBar = value
	}
}

// WithColors set hide ratio
func WithColors(value *Colors) Option {
	return func(w *Widget) {
		w.Colors = value
	}
}
