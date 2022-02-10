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
	updateChan chan string
	Output     io.Writer
	active     bool
	closed     bool
	style      *gchalk.Builder
	MaxWidth   int
	HideCursor bool
}

//New creates a new Spinner instance
func New(options ...Option) *Widget {
	s := &Widget{
		Label:      "Executing...",
		Frames:     FramesBall,
		FrameRate:  frameRate,
		style:      gchalk.WithBrightCyan(),
		runChan:    make(chan struct{}),
		updateChan: make(chan string),
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

//Close will stop taking updates
func (s *Widget) Close() {
	s.Lock()
	defer s.Unlock()
	if !s.closed {
		close(s.updateChan)
		s.closed = true
	}
}

//UpdateLabel with update text to be rendered in the screen
func (s *Widget) UpdateLabel(l string) {
	if s.active {
		s.updateChan <- l
	}
}

//SetStyle will set the style used in the label
func (s *Widget) SetStyle(style *gchalk.Builder) *Widget {
	s.Lock()
	s.style = style
	s.Unlock()
	return s
}

//ApplyStyle will apply the style to the given string
func (s *Widget) ApplyStyle(str string) string {
	if s.style != nil {
		str = s.style.Paint(str)
	}
	return str
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

	var l string
	var ok bool

	for {
		select {
		case <-s.runChan:
			return
		case l, ok = <-s.updateChan:
			if ok {
				s.SetLabel(l)
			}
		default:
			s.animate()
		}
	}
}

func (s *Widget) animate() {
	var out string

	//hides cursor
	if s.HideCursor {
		fmt.Fprint(s.Output, "\033[?25l")
	}

	for i := 0; i < len(s.Frames); i++ {
		frame := s.Frames[i]

		frame = s.ApplyStyle(frame)

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
