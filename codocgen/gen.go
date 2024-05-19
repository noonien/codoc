package codocgen

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"

	"github.com/noonien/codoc"
	"golang.org/x/tools/go/packages"
)

// RegisterPath registers a package at the given path.
func RegisterPath(path string, opts ...Option) error {
	pkg, err := FromPath(path, opts...)
	if err != nil {
		return err
	}

	codoc.Register(*pkg)
	return nil
}

// FromPath generates documentation for a package at the given path.
func FromPath(path string, opts ...Option) (*codoc.Package, error) {
	conf := &config{}
	for _, opt := range opts {
		opt(conf)
	}

	info, err := getInfo(path)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse package %q: %v", path, err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages in %q", path)
	}
	if len(pkgs) > 1 {
		return nil, fmt.Errorf("multiple packages in %q", path)
	}

	pkgast := pkgs[info.Name]
	pkgdoc := doc.New(pkgast, info.ID, doc.AllDecls)

	funcs := make(map[string]codoc.Function, len(pkgdoc.Funcs))
	for _, fn := range pkgdoc.Funcs {
		fn := getFunc(fn)
		if conf.filterFunc(fn) {
			funcs[fn.Name] = fn
		}
	}

	structs := make(map[string]codoc.Struct, len(pkgdoc.Types))
	for _, typ := range pkgdoc.Types {
		ts := typ.Decl.Specs[0].(*ast.TypeSpec)
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			continue
		}

		for _, fn := range typ.Funcs {
			fn := getFunc(fn)
			if conf.filterFunc(fn) {
				funcs[fn.Name] = fn
			}
		}

		methods := make(map[string]codoc.Function, len(typ.Methods))
		for _, fn := range typ.Methods {
			m := getFunc(fn)
			if conf.filterFunc(m) {
				methods[m.Name] = m
			}
		}

		fields := map[string]codoc.Field{}
		for _, field := range st.Fields.List {
			for _, name := range field.Names {
				doc := strings.TrimSpace(field.Doc.Text())
				comment := strings.TrimSpace(field.Comment.Text())
				if len(doc) > 0 || len(comment) > 0 {
					fields[name.Name] = codoc.Field{
						Doc:     doc,
						Comment: comment,
					}
				}
			}
		}

		cst := codoc.Struct{
			Name:    typ.Name,
			Doc:     strings.TrimSpace(typ.Doc),
			Fields:  fields,
			Methods: methods,
		}

		if conf.filterStruct(cst) {
			structs[typ.Name] = cst
		}
	}

	return &codoc.Package{
		Name:      info.Name,
		ID:        info.ID,
		Doc:       strings.TrimSpace(pkgdoc.Doc),
		Functions: funcs,
		Structs:   structs,
	}, nil
}

type PackageError []packages.Error

func (PackageError) Error() string { return "package contains errors" }

func getInfo(path string) (*packages.Package, error) {
	infos, err := packages.Load(nil, path)
	if err != nil {
		return nil, fmt.Errorf("load package %q: %v", path, err)
	}

	if len(infos) == 0 {
		return nil, fmt.Errorf("no packages in %q", path)
	}
	if len(infos) > 1 {
		return nil, fmt.Errorf("multiple packages in %q", path)
	}

	info := infos[0]
	if len(info.Errors) > 0 {
		return nil, PackageError(info.Errors)
	}

	return info, nil
}

func getFunc(fn *doc.Func) codoc.Function {
	dt := fn.Decl.Type

	var args []string
	if dt.Params != nil {
		for _, arg := range dt.Params.List {
			for _, ident := range arg.Names {
				if len(ident.Name) > 0 {
					args = append(args, ident.Name)
				}
			}
		}
	}

	var results []string
	if dt.Results != nil {
		for _, res := range dt.Results.List {
			for _, ident := range res.Names {
				if len(ident.Name) > 0 {
					results = append(results, ident.Name)
				}
			}
		}
	}

	return codoc.Function{
		Name:    fn.Name,
		Doc:     strings.TrimSpace(fn.Doc),
		Args:    args,
		Results: results,
	}
}
