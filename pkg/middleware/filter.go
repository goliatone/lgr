package middleware

import (
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"github.com/dop251/goja"
	"github.com/goliatone/lgr/pkg/logging"
	"github.com/goliatone/lgr/pkg/render"
)

// Middleware is the function used to run a line handler
type Middleware interface {
	Next(msg *logging.Message, opts *render.Options)
}

type filter struct {
	script string
}

func (f *filter) Next(msg *logging.Message, _ *render.Options) bool {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	vm.Set("line", msg)

	vm.Set("addField", func(key string, value interface{}) bool {
		msg.AddField(key, value)
		return true
	})

	vm.Set("deleteFields", func(keys ...string) bool {
		msg.DeleteFields(keys...)
		return true
	})

	val, err := vm.RunString(f.script)
	if err != nil {
		return true
	}

	if next, ok := val.Export().(bool); ok {
		return next
	}

	return true
}

type FilterSet struct {
	filters []filter
}

func (f *FilterSet) Next(msg *logging.Message, opts *render.Options) bool {
	for _, f := range f.filters {
		if !f.Next(msg, opts) {
			return false
		}
	}
	return true
}

func NewFilterSet(filters []string) (FilterSet, error) {
	o := []filter{}
	for _, f := range filters {
		if strings.HasPrefix(f, "@") {
			file, err := ioutil.ReadFile(trimFirstRune(f))
			if err != nil {
				return FilterSet{}, err
			}
			f = string(file)
		}

		o = append(o, filter{script: f})
	}

	return FilterSet{
		filters: o,
	}, nil
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
