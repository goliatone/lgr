package widgets

//TODO: move to widgets/spinner package
import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/jwalton/gchalk"
)

const spinnerFrameRate = time.Millisecond * 100

//SpinnerFrames ...
type SpinnerFrames []string

// var defaultFrames = SpinnerFrames{"|", "/", "-", "\\"}
// var defaultFrames = SpinnerFrames{"⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽", "⣾"}
// var defaultFrames = SpinnerFrames{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"}
// var defaultFrames = SpinnerFrames{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
// var defaultFrames = SpinnerFrames{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"}

var (
	FramesBall          = SpinnerFrames{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"}
	FramesBarHorizontal = SpinnerFrames{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "▊", "▋", "▌", "▍", "▎"}
	FramesBarsVertical  = SpinnerFrames{"▁", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃", "▁"}
)

//Option will apply the provided configuration
//to a spinner instance.
type Option func(*Spinner)

//WithStyle adds color style
func WithStyle(style *gchalk.Builder) Option {
	return func(s *Spinner) {
		s.Style = style
	}
}

//Spinner is the spinner struct
type Spinner struct {
	sync.Mutex
	Label      string
	Frames     SpinnerFrames
	FrameRate  time.Duration
	runChan    chan struct{}
	Output     io.Writer
	active     bool
	Style      *gchalk.Builder
	HideCursor bool
}

//NewSpinner creates a new spinner...
func NewSpinner(label string, options ...Option) *Spinner {
	s := &Spinner{
		Label:     label,
		Frames:    FramesBall,
		FrameRate: spinnerFrameRate,
		Style:     gchalk.WithBrightCyan(),
		runChan:   make(chan struct{}),
	}

	for _, option := range options {
		option(s)
	}

	return s
}

//Start starts the spinner
func (s *Spinner) Start() *Spinner {
	s.Lock()

	if s.active {
		s.Unlock()
		return s
	}

	s.active = true
	s.Unlock()

	go s.writter()
	return s
}

//Stop will stop the spinner
func (s *Spinner) Stop() *Spinner {
	s.Lock()
	defer s.Unlock()

	if !s.active {
		return s
	}

	s.active = false
	close(s.runChan)

	if s.HideCursor {
		//show cursor again
		fmt.Fprint(s.Output, "\033[?25h")
	}

	s.clearOutput()

	return s
}

//SetLabel updates the label value
func (s *Spinner) SetLabel(label string) *Spinner {
	s.Lock()
	s.Label = label
	s.Unlock()
	return s
}

func (s *Spinner) clearOutput() {
	fmt.Fprintf(s.Output, "\r\033[K")
}

func (s *Spinner) writter() {
	s.animate()
	for {
		select {
		case <-s.runChan:
			return
		default:
			s.animate()
		}
	}
}

func (s *Spinner) animate() {
	var out string

	style := gchalk.WithBrightCyan()

	//hides cursor
	if s.HideCursor {
		fmt.Fprint(s.Output, "\033[?25l")
	}

	for i := 0; i < len(s.Frames); i++ {
		frame := s.Frames[i]

		if s.Style != nil {
			frame = style.Paint(frame)
		}

		out = fmt.Sprintf("\r%s %s", frame, s.Label)

		fmt.Fprint(s.Output, out)

		time.Sleep(s.FrameRate)
		s.clearOutput()
	}
}
