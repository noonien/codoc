package testpkg

// ExportedFunc is an exported function
func ExportedFunc() {}

// unexportedFunc is an unexported function
func unexportedFunc() {}

// ExportedType is an exported struct
type ExportedType struct{}

// unexportedType is an unexported struct
type unexportedType struct{}
