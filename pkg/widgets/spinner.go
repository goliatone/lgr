package widgets

import (
	"fmt"
	"io"
	"sync"
	"time"
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

//Spinner is the spinner struct
type Spinner struct {
	sync.Mutex
	Label     string
	Frames    SpinnerFrames
	FrameRate time.Duration
	runChan   chan struct{}
	Output    io.Writer
	active    bool
}

//NewSpinner creates a new spinner...
func NewSpinner(label string) *Spinner {
	return &Spinner{
		Label:     label,
		Frames:    FramesBarHorizontal,
		FrameRate: spinnerFrameRate,
		runChan:   make(chan struct{}),
	}
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

	for i := 0; i < len(s.Frames); i++ {
		out = fmt.Sprintf("\r%s %s", s.Frames[i], s.Label)
		fmt.Fprint(s.Output, out)
		time.Sleep(s.FrameRate)
		s.clearOutput()
	}
}
