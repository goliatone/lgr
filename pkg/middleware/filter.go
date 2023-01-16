package middleware

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"github.com/dop251/goja"
	"github.com/goliatone/lgr/pkg/logging"
	"github.com/goliatone/lgr/pkg/render"
)

type Filter struct {
	Script string
}

func (f *Filter) Next(msg *logging.Message, _ *render.Options) bool {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	vm.Set("line", msg)

	vm.Set("addField", func(key string, value string) bool {
		msg.Fields = append(msg.Fields, &logging.MessageField{
			Key:   key,
			Value: value,
		})
		return true
	})

	val, err := vm.RunString(f.Script)
	if err != nil {
		return true
	}

	if next, ok := val.Export().(bool); ok {
		return next
	}

	return true
}

type FilterSet struct {
	filters []Filter
}

func (f *FilterSet) Next(msg *logging.Message, opts *render.Options) bool {
	for _, f := range f.filters {
		if !f.Next(msg, opts) {
			return false
		}
	}
	return true
}

func NewFilterSet(filters []string) FilterSet {
	o := []Filter{}
	for _, f := range filters {
		fmt.Printf("add script: %s\n", f)
		if strings.HasPrefix(f, "@") {
			script, err := ioutil.ReadFile(trimFirstRune(f))
			if err != nil {
				continue
			}
			o = append(o, Filter{Script: string(script)})
		} else {
			o = append(o, Filter{Script: f})
		}
	}

	return FilterSet{
		filters: o,
	}
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
