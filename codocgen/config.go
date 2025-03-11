// Package codocgen provides tools for generating code documentation from Go source code.
// This package analyzes Go packages and produces codoc.Package objects with documentation.
package codocgen

import (
	"unicode"
	"unicode/utf8"

	"github.com/noonien/codoc"
)

// Option is a function type that modifies a config.
// Options are used to customize the behavior of the documentation generator.
type Option func(*config)

// config holds the configuration for the documentation generator.
// It contains filters for functions and structs to determine what gets included in the documentation.
type config struct {
	funcFilter   []func(fn codoc.Function) bool // Filters for functions
	structFilter []func(st codoc.Struct) bool   // Filters for structs
}

// FilterFuncs adds a function filter to the configuration.
// The filter function takes a Function and returns true if it should be included in the documentation.
func FilterFuncs(fn func(fn codoc.Function) bool) Option {
	return func(c *config) {
		c.funcFilter = append(c.funcFilter, fn)
	}
}

// FilterStructs adds a struct filter to the configuration.
// The filter function takes a Struct and returns true if it should be included in the documentation.
func FilterStructs(fn func(st codoc.Struct) bool) Option {
	return func(c *config) {
		c.structFilter = append(c.structFilter, fn)
	}
}

// Exported returns an Option that filters to include only exported functions and structs.
// Exported items are those that start with an uppercase letter.
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

// WithDoc returns an Option that filters to include only functions and structs with documentation.
// This is useful to ensure that only documented code appears in the output.
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

// filterFunc applies all function filters in the configuration to a function.
// Returns true only if all filters return true, meaning the function should be included.
func (c *config) filterFunc(fn codoc.Function) bool {
	for _, f := range c.funcFilter {
		if !f(fn) {
			return false
		}
	}
	return true
}

// filterStruct applies all struct filters in the configuration to a struct.
// Returns true only if all filters return true, meaning the struct should be included.
func (c *config) filterStruct(st codoc.Struct) bool {
	for _, f := range c.structFilter {
		if !f(st) {
			return false
		}
	}
	return true
}
