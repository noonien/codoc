package codoc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterAndGetPackage(t *testing.T) {
	// Create a simple package
	testPkg := Package{
		ID:   "example.com/testpkg",
		Name: "testpkg",
		Doc:  "This is a test package",
		Functions: map[string]Function{
			"TestFunc": {
				Name:    "TestFunc",
				Doc:     "Test function documentation",
				Args:    []string{"arg1", "arg2"},
				Results: []string{"result1", "result2"},
			},
		},
		Structs: map[string]Struct{
			"TestStruct": {
				Name: "TestStruct",
				Doc:  "Test struct documentation",
				Fields: map[string]Field{
					"Field1": {
						Name:    "Field1",
						Doc:     "Field1 documentation",
						Comment: "Field1 comment",
					},
				},
				Methods: map[string]Function{
					"Method1": {
						Name:    "Method1",
						Doc:     "Method1 documentation",
						Args:    []string{"arg1"},
						Results: []string{"result1"},
					},
				},
			},
		},
	}

	// Register the package
	Register(testPkg)

	// Test GetPackage
	pkg := GetPackage("example.com/testpkg")
	require.NotNil(t, pkg, "GetPackage returned nil for registered package")
	assert.Equal(t, "example.com/testpkg", pkg.ID, "Package ID mismatch")
	assert.Equal(t, "testpkg", pkg.Name, "Package name mismatch")
	assert.Equal(t, "This is a test package", pkg.Doc, "Package doc mismatch")

	// Test GetFunction
	fn := GetFunction("example.com/testpkg.TestFunc")
	require.NotNil(t, fn, "GetFunction returned nil for registered function")
	assert.Equal(t, "TestFunc", fn.Name, "Function name mismatch")
	assert.Equal(t, "Test function documentation", fn.Doc, "Function doc mismatch")
	assert.Equal(t, []string{"arg1", "arg2"}, fn.Args, "Function args mismatch")
	assert.Equal(t, []string{"result1", "result2"}, fn.Results, "Function results mismatch")

	// Test GetStruct
	st := GetStruct("example.com/testpkg.TestStruct")
	require.NotNil(t, st, "GetStruct returned nil for registered struct")
	assert.Equal(t, "TestStruct", st.Name, "Struct name mismatch")
	assert.Equal(t, "Test struct documentation", st.Doc, "Struct doc mismatch")

	// Test Field
	field, ok := st.Fields["Field1"]
	assert.True(t, ok, "Field 'Field1' not found in struct")
	assert.Equal(t, "Field1", field.Name, "Field name mismatch")
	assert.Equal(t, "Field1 documentation", field.Doc, "Field doc mismatch")
	assert.Equal(t, "Field1 comment", field.Comment, "Field comment mismatch")

	// Test Method
	method, ok := st.Methods["Method1"]
	assert.True(t, ok, "Method 'Method1' not found in struct")
	assert.Equal(t, "Method1", method.Name, "Method name mismatch")
	assert.Equal(t, "Method1 documentation", method.Doc, "Method doc mismatch")

	methodfn := GetFunction("example.com/testpkg.TestStruct.Method1")
	assert.NotNil(t, methodfn, "GetFunction returned nil for registered method")
	assert.Equal(t, "Method1", methodfn.Name, "Method name mismatch")
	assert.Equal(t, "Method1 documentation", methodfn.Doc, "Method doc mismatch")
}

func TestGetNonExistentItems(t *testing.T) {
	// Test getting a package that doesn't exist
	pkg := GetPackage("nonexistent.pkg")
	assert.Nil(t, pkg, "GetPackage should return nil for non-existent package")

	// Test getting a function that doesn't exist
	fn := GetFunction("nonexistent.pkg.SomeFunc")
	assert.Nil(t, fn, "GetFunction should return nil for non-existent function")

	// Test getting a struct that doesn't exist
	st := GetStruct("nonexistent.pkg.SomeStruct")
	assert.Nil(t, st, "GetStruct should return nil for non-existent struct")
}

func TestRegisterMainPackage(t *testing.T) {
	// Create a simple "main" package
	mainPkg := Package{
		ID:   "example.com/main",
		Name: "main",
		Doc:  "This is the main package",
		Functions: map[string]Function{
			"MainFunc": {
				Name: "MainFunc",
				Doc:  "Main function documentation",
			},
		},
	}

	// Register the main package
	Register(mainPkg)

	// Main packages should be registered with the ID "main"
	pkg := GetPackage("main")
	require.NotNil(t, pkg, "GetPackage returned nil for registered main package")
	assert.Equal(t, "main", pkg.Name, "Package name mismatch")
	assert.Equal(t, "This is the main package", pkg.Doc, "Package doc mismatch")

	// Test the main package function
	fn := GetFunction("main.MainFunc")
	require.NotNil(t, fn, "GetFunction returned nil for registered main function")
	assert.Equal(t, "MainFunc", fn.Name, "Function name mismatch")
}
