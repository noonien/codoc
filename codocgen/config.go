package codocgen

import (
	"unicode"
	"unicode/utf8"

	"github.com/noonien/codoc"
)

type Option func(*config)

type config struct {
	funcFilter   []func(fn codoc.Function) bool
	structFilter []func(st codoc.Struct) bool
}

func FilterFuncs(fn func(fn codoc.Function) bool) Option {
	return func(c *config) {
		c.funcFilter = append(c.funcFilter, fn)
	}
}

func FilterStructs(fn func(st codoc.Struct) bool) Option {
	return func(c *config) {
		c.structFilter = append(c.structFilter, fn)
	}
}

func Exported() Option {
	return func(c *config) {
		c.funcFilter = append(c.funcFilter, func(fn codoc.Function) bool {
			r, _ := utf8.DecodeRuneInString(fn.Name)
			return unicode.IsUpper(r)
		})

		c.structFilter = append(c.structFilter, func(st codoc.Struct) bool {
			r, _ := utf8.DecodeRuneInString(st.Name)
			return unicode.IsUpper(r)
		})
	}
}

func WithDoc() Option {
	return func(c *config) {
		c.funcFilter = append(c.funcFilter, func(fn codoc.Function) bool {
			return fn.Doc != ""
		})

		c.structFilter = append(c.structFilter, func(st codoc.Struct) bool {
			return st.Doc != ""
		})
	}
}

func (c *config) filterFunc(fn codoc.Function) bool {
	for _, f := range c.funcFilter {
		if !f(fn) {
			return false
		}
	}
	return true
}

func (c *config) filterStruct(st codoc.Struct) bool {
	for _, f := range c.structFilter {
		if !f(st) {
			return false
		}
	}
	return true
}
