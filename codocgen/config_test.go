package codocgen

import (
	"testing"

	"github.com/noonien/codoc"
	"github.com/stretchr/testify/assert"
)

func TestFilterFuncs(t *testing.T) {
	c := &config{}

	// Add a filter that only accepts functions named "AcceptedFunc"
	FilterFuncs(func(fn codoc.Function) bool {
		return fn.Name == "AcceptedFunc"
	})(c)

	// Test with an accepted function
	assert.True(t, c.filterFunc(codoc.Function{Name: "AcceptedFunc"}), "Function 'AcceptedFunc' should be accepted")

	// Test with a non-accepted function
	assert.False(t, c.filterFunc(codoc.Function{Name: "RejectedFunc"}), "Function 'RejectedFunc' should be rejected")
}

func TestFilterStructs(t *testing.T) {
	c := &config{}

	// Add a filter that only accepts structs named "AcceptedStruct"
	FilterStructs(func(st codoc.Struct) bool {
		return st.Name == "AcceptedStruct"
	})(c)

	// Test with an accepted struct
	assert.True(t, c.filterStruct(codoc.Struct{Name: "AcceptedStruct"}), "Struct 'AcceptedStruct' should be accepted")

	// Test with a non-accepted struct
	assert.False(t, c.filterStruct(codoc.Struct{Name: "RejectedStruct"}), "Struct 'RejectedStruct' should be rejected")
}

func TestExported(t *testing.T) {
	c := &config{}

	// Add exported filter
	Exported()(c)

	// Test with exported function
	assert.True(t, c.filterFunc(codoc.Function{Name: "ExportedFunc"}), "Function 'ExportedFunc' should be accepted")

	// Test with unexported function
	assert.False(t, c.filterFunc(codoc.Function{Name: "unexportedFunc"}), "Function 'unexportedFunc' should be rejected")

	// Test with exported struct
	assert.True(t, c.filterStruct(codoc.Struct{Name: "ExportedStruct"}), "Struct 'ExportedStruct' should be accepted")

	// Test with unexported struct
	assert.False(t, c.filterStruct(codoc.Struct{Name: "unexportedStruct"}), "Struct 'unexportedStruct' should be rejected")
}

func TestWithDoc(t *testing.T) {
	c := &config{}

	// Add WithDoc filter
	WithDoc()(c)

	// Test with documented function
	assert.True(t, c.filterFunc(codoc.Function{Name: "DocFunc", Doc: "This function has docs"}),
		"Function with documentation should be accepted")

	// Test with undocumented function
	assert.False(t, c.filterFunc(codoc.Function{Name: "NoDocFunc", Doc: ""}),
		"Function without documentation should be rejected")

	// Test with documented struct
	assert.True(t, c.filterStruct(codoc.Struct{Name: "DocStruct", Doc: "This struct has docs"}),
		"Struct with documentation should be accepted")

	// Test with undocumented struct
	assert.False(t, c.filterStruct(codoc.Struct{Name: "NoDocStruct", Doc: ""}),
		"Struct without documentation should be rejected")
}

func TestMultipleFilters(t *testing.T) {
	c := &config{}

	// Add multiple filters
	Exported()(c)
	WithDoc()(c)

	// Test with exported and documented function
	assert.True(t, c.filterFunc(codoc.Function{Name: "ExportedDocFunc", Doc: "This function has docs"}),
		"Exported function with documentation should be accepted")

	// Test with exported but undocumented function
	assert.False(t, c.filterFunc(codoc.Function{Name: "ExportedNoDocFunc", Doc: ""}),
		"Exported function without documentation should be rejected")

	// Test with unexported but documented function
	assert.False(t, c.filterFunc(codoc.Function{Name: "unexportedDocFunc", Doc: "This function has docs"}),
		"Unexported function with documentation should be rejected")

	// Similar tests for structs
	assert.True(t, c.filterStruct(codoc.Struct{Name: "ExportedDocStruct", Doc: "This struct has docs"}),
		"Exported struct with documentation should be accepted")

	assert.False(t, c.filterStruct(codoc.Struct{Name: "ExportedNoDocStruct", Doc: ""}),
		"Exported struct without documentation should be rejected")

	assert.False(t, c.filterStruct(codoc.Struct{Name: "unexportedDocStruct", Doc: "This struct has docs"}),
		"Unexported struct with documentation should be rejected")
}
