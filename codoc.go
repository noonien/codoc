package codoc

import "sync"

type Package struct {
	ID        string
	Name      string
	Doc       string
	Functions map[string]Function
	Structs   map[string]Struct
}

type Function struct {
	Name    string
	Doc     string
	Args    []string
	Results []string
}

type Struct struct {
	Name    string
	Doc     string
	Fields  map[string]Field
	Methods map[string]Function
}

type Field struct {
	Name    string
	Doc     string
	Comment string
}

var funcs = map[string]Function{}
var strucsts = map[string]Struct{}
var pkgs = map[string]Package{}
var mu sync.RWMutex

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

func GetPackage(id string) *Package {
	mu.RLock()
	defer mu.RUnlock()
	pkg, ok := pkgs[id]
	if !ok {
		return nil
	}
	return &pkg
}

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

func GetStruct(id string) *Struct {
	mu.RLock()
	defer mu.RUnlock()

	st, ok := strucsts[id]
	if !ok {
		return nil
	}
	return &st
}
