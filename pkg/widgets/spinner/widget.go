package spinner

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/jwalton/gchalk"
)

const frameRate = time.Millisecond * 100

//Frames ...
type Frames []string

//Widget is the spinner struct
type Widget struct {
	sync.Mutex
	Label      string
	Frames     Frames
	FrameRate  time.Duration
	runChan    chan struct{}
	Output     io.Writer
	active     bool
	Style      *gchalk.Builder
	MaxWidth   int
	HideCursor bool
}

//New creates a new Spinner instance
func New(options ...Option) *Widget {
	s := &Widget{
		Label:     "Executing...",
		Frames:    FramesBall,
		FrameRate: frameRate,
		Style:     gchalk.WithBrightCyan(),
		runChan:   make(chan struct{}),
	}

	for _, option := range options {
		option(s)
	}

	return s
}

//Start starts the spinner
func (s *Widget) Start() *Widget {
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
func (s *Widget) Stop() *Widget {
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
func (s *Widget) SetLabel(label string) *Widget {
	s.Lock()
	s.Label = label
	s.Unlock()
	return s
}

func (s *Widget) clearOutput() {
	fmt.Fprintf(s.Output, "\r\033[K")
}

func (s *Widget) writter() {
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

func (s *Widget) animate() {
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

		label := clipString(s.Label, s.MaxWidth)
		out = fmt.Sprintf("\r %s %s", frame, label)

		fmt.Fprint(s.Output, out)

		time.Sleep(s.FrameRate)
		s.clearOutput()
	}
}

func clipString(str string, width int) string {
	if len(str) < width {
		return str
	}
	return str[0:width-4] + "..."
}
