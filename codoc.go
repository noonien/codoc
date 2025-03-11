// Package codoc provides functionality for storing and retrieving code documentation.
// It defines data structures for representing code elements like packages, functions, and structs,
// and provides an API for registering and retrieving documentation.
package codoc

import (
	"strings"
	"sync"
)

// Package represents a Go package with its documentation.
// It contains information about the package itself, as well as
// maps of the functions and structs defined within it.
type Package struct {
	ID        string              // Unique identifier for the package
	Name      string              // Package name
	Doc       string              // Package documentation string
	Functions map[string]Function // Map of functions in the package
	Structs   map[string]Struct   // Map of structs in the package
}

// Function represents a Go function with its documentation.
// It includes the function's name, documentation, and parameter information.
type Function struct {
	Name    string   // Function name
	Doc     string   // Function documentation string
	Args    []string // List of argument names
	Results []string // List of result names
}

// Struct represents a Go struct with its documentation.
// It includes the struct's name, documentation, fields, and methods.
type Struct struct {
	Name    string              // Struct name
	Doc     string              // Struct documentation string
	Fields  map[string]Field    // Map of fields in the struct
	Methods map[string]Function // Map of methods associated with the struct
}

// Field represents a field in a struct with its documentation.
type Field struct {
	Name    string // Field name
	Doc     string // Field documentation string
	Comment string // Inline comment for the field
}

// Global maps to store registered functions, structs, and packages
var funcs = map[string]Function{}
var strucsts = map[string]Struct{}
var pkgs = map[string]Package{}
var mu sync.RWMutex // Mutex to protect concurrent access to the maps

// Register adds a package and all its components to the global registry.
// It uses the package's ID as a prefix for registering functions and structs.
func Register(pkg Package) {
	mu.Lock()
	defer mu.Unlock()

	id := pkg.ID
	if pkg.Name == "main" {
		id = "main"
	}
	pkgs[id] = pkg
	prefix := id + "."
	for _, fn := range pkg.Functions {
		funcs[prefix+fn.Name] = fn
	}
	for _, st := range pkg.Structs {
		strucsts[prefix+st.Name] = st
	}
}

// GetPackage retrieves a package from the registry by its ID.
// Returns nil if the package is not found.
func GetPackage(id string) *Package {
	mu.RLock()
	defer mu.RUnlock()
	pkg, ok := pkgs[id]
	if !ok {
		return nil
	}
	return &pkg
}

// GetFunction retrieves a function from the registry by its ID.
// The ID can be either a direct function ID or a struct method ID (pkg.struct.method).
// Returns nil if the function is not found.
func GetFunction(id string) *Function {
	mu.RLock()
	defer mu.RUnlock()

	fn, ok := funcs[id]
	if ok {
		return &fn
	}
	// Find the last dot in the id to split into struct and method parts
	lastDotIndex := strings.LastIndex(id, ".")
	if lastDotIndex == -1 {
		return nil
	}

	// Extract the struct and method names
	structID := id[:lastDotIndex]
	methodname := id[lastDotIndex+1:]

	// Get the struct
	st := GetStruct(structID)
	if st == nil {
		return nil
	}

	// Get the method from the struct
	fn, exists := st.Methods[methodname]
	if !exists {
		return nil
	}

	return &fn
}

// GetStruct retrieves a struct from the registry by its ID.
// Returns nil if the struct is not found.
func GetStruct(id string) *Struct {
	mu.RLock()
	defer mu.RUnlock()

	st, ok := strucsts[id]
	if !ok {
		return nil
	}
	return &st
}
