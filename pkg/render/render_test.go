package render

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/goliatone/lgr/pkg/logging"
)

func TestPrint(t *testing.T) {

	type args struct {
		msg  logging.Message
		opts Options
	}

	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				msg: logging.Message{
					Message:   "this is a test",
					Level:     "debug",
					Timestamp: nil,
				},
				opts: Options{
					Level:           "debug",
					Color:           "neutral",
					Heading:         "H",
					Bold:            false,
					TimestampFormat: TimestampFormat,
					Writer:          os.Stdout,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Print(&tt.args.msg, &tt.args.opts)
		})
	}
}

func Test_checkInput(t *testing.T) {

	f, df, err := mockStdin(t, "hello world")
	if err != nil {
		t.Fatal(err)
	}
	defer df()

	has, err := checkInput(f)
	if !has {
		t.Error("checkInput() should detect content in stdin")
	}

	f2, df2, err := mockStdin(t, "")
	if err != nil {
		t.Fatal(err)
	}
	defer df2()

	has, err = checkInput(f2)
	if has {
		t.Error("checkInput() should detect content in stdin")
	}
}

func Test_streamToString(t *testing.T) {
	type args struct {
		s io.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "content",
			args: args{s: strings.NewReader("expected")},
			want: "expected",
		},
		{
			name: "empty",
			args: args{s: strings.NewReader("")},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := streamToString(tt.args.s); got != tt.want {
				t.Errorf("streamToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeStyles(t *testing.T) {
	type args struct {
		styles []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "single",
			args: args{[]string{"hi-black"}},
			want: []string{"brightBlack"},
		},
		{
			name: "multiple",
			args: args{[]string{"hi-black", "bg-black", "hi-bg-magenta"}},
			want: []string{"brightBlack", "bgBlack", "bgBrightMagenta"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeStyles(tt.args.styles...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalizeStyles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockStdin(t *testing.T, input string) (*os.File, func(), error) {
	t.Helper()

	tmpf, err := ioutil.TempFile(t.TempDir(), t.Name())
	if err != nil {
		return nil, nil, err
	}
	content := []byte(input)

	if _, err := tmpf.Write(content); err != nil {
		return nil, nil, err
	}

	if _, err := tmpf.Seek(0, 0); err != nil {
		return nil, nil, err
	}

	return tmpf, func() {
		os.Remove(tmpf.Name())
	}, nil
}
