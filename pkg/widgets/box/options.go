package box

//Option will apply the provided configuration
//to  a spinner instance
type Option func(*Widget)

//WithWidth configures width option
func WithWidth(value int) Option {
	return func(w *Widget) {
		w.Width = value
	}
}

//WithMaxWidth configures max width option
func WithMaxWidth(value int) Option {
	return func(w *Widget) {
		w.MaxWidth = value
	}
}

//WithHeight configures max width option
func WithHeight(value int) Option {
	return func(w *Widget) {
		w.Height = value
	}
}

//WithTitle configures title option
func WithTitle(value string) Option {
	return func(w *Widget) {
		w.Title = value
	}
}

//WithContent configures content option
func WithContent(value string) Option {
	return func(w *Widget) {
		w.Content = value
	}
}

//WithTemplate configures content option
func WithTemplate(value string) Option {
	return func(w *Widget) {
		w.Template = value
	}
}

//WithHpad configures max width option
func WithHpad(value int) Option {
	return func(w *Widget) {
		w.Hpad = value
	}
}

//WithVpad configures max width option
func WithVpad(value int) Option {
	return func(w *Widget) {
		w.Vpad = value
	}
}

//WithScreenW configures max width option
func WithScreenW(value int) Option {
	return func(w *Widget) {
		w.ScreenW = value
	}
}

//WithScreenH configures max width option
func WithScreenH(value int) Option {
	return func(w *Widget) {
		w.ScreenH = value
	}
}

//WithAlignment configures max width option
func WithAlignment(value string) Option {
	return func(w *Widget) {
		w.Alignment = value
	}
}
