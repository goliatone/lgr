package middleware

import (
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

func NewFilterSet(filters []string) (FilterSet, error) {
	o := []Filter{}
	for _, filter := range filters {
		if strings.HasPrefix(filter, "@") {
			file, err := ioutil.ReadFile(trimFirstRune(filter))
			if err != nil {
				return FilterSet{}, err
			}
			filter = string(file)
		}

		o = append(o, Filter{Script: filter})
	}

	return FilterSet{
		filters: o,
	}, nil
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
