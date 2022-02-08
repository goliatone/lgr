package spinner

import (
	"io"
	"time"

	"github.com/jwalton/gchalk"
)

//Option will apply the provided configuration
//to a spinner instance.
type Option func(*Widget)

//WithLabel adds color style
func WithLabel(value string) Option {
	return func(s *Widget) {
		s.Label = value
	}
}

//WithFrames adds color style
func WithFrames(value Frames) Option {
	return func(s *Widget) {
		s.Frames = value
	}
}

//WithFrameRate adds color style
func WithFrameRate(value time.Duration) Option {
	return func(s *Widget) {
		s.FrameRate = value
	}
}

//WithStyle adds color style
func WithStyle(value *gchalk.Builder) Option {
	return func(s *Widget) {
		s.Style = value
	}
}

//WithOutput adds color style
func WithOutput(value io.Writer) Option {
	return func(s *Widget) {
		s.Output = value
	}
}

//WithMaxWidth adds color style
func WithMaxWidth(value int) Option {
	return func(s *Widget) {
		s.MaxWidth = value
	}
}

//WithHideCursor adds color style
func WithHideCursor(value bool) Option {
	return func(s *Widget) {
		s.HideCursor = value
	}
}
