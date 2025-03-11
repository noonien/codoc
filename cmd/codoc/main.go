// Package main provides a command-line tool for generating code documentation.
// The tool analyzes Go packages and generates code that can be used to register
// documentation information with the codoc package.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/alecthomas/repr"
	"github.com/noonien/codoc"
	"github.com/noonien/codoc/codocgen"
)

// Command-line flags
var (
	outFile  = flag.String("out", "", "output file, leave empty to write to stdout")
	pkgName  = flag.String("pkg", "", "output file package")
	exported = flag.Bool("e", false, "only register exported functions and structs")
)

// main is the entry point for the codoc command-line tool.
// It parses command-line flags, processes the specified packages,
// and generates documentation in the desired output format.
func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	// Parse command-line flags
	flag.Parse()
	if len(*pkgName) == 0 {
		flag.Usage()
		log.Fatal("missing flag: pkg")
	}

	// Check for package paths
	paths := flag.Args()
	if len(paths) == 0 {
		flag.Usage()
		log.Fatalf("no package paths specified")
	}

	// Set up documentation generation options
	opts := []codocgen.Option{}
	if *exported {
		opts = append(opts, codocgen.Exported())
	}

	// Process each package and extract documentation
	var pkgs []*codoc.Package
	for _, p := range flag.Args() {
		pkg, err := codocgen.FromPath(p, opts...)
		if err != nil {
			log.Fatalf("could not get docs for %q: %v", p, err)
		}
		log.Printf("got docs for %s", pkg.Name)
		pkgs = append(pkgs, pkg)
	}

	// Set up output file
	var f *os.File
	if *outFile == "" || *outFile == "-" {
		f = os.Stdout
	} else {
		var err error
		f, err = os.Create(*outFile)
		if err != nil {
			log.Fatalf("cannot create file: %v", err)
		}
		defer f.Close()
	}

	// Set up gofmt to format the output
	gofmt := exec.Command("gofmt", "-s")

	fmtw, err := gofmt.StdinPipe()
	if err != nil {
		log.Fatalf("cannot get stdin pipe: %v", err)
	}
	gofmt.Stdout = f
	gofmt.Stderr = os.Stderr

	if err := gofmt.Start(); err != nil {
		log.Fatalf("cannot start gofmt: %v", err)
	}
	writeDoc(fmtw, pkgs)
	if err := gofmt.Wait(); err != nil {
		log.Fatal(err)
	}
}

// writeDoc generates the Go code to register documentation for packages.
// It writes the code to the specified writer, which is piped through gofmt.
// The generated code includes imports and a call to codoc.Register for each package.
func writeDoc(w io.WriteCloser, pkgs []*codoc.Package) {
	defer w.Close()

	// Write file header with timestamp
	fmt.Fprintf(w, "// generated @ %s by gendoc\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(w, "package %s\n", *pkgName)
	fmt.Fprintln(w)
	io.WriteString(w, "import \"github.com/noonien/codoc\"\n")
	fmt.Fprintln(w)

	// Write init function that registers all packages
	io.WriteString(w, "func init() {\n")
	for _, pkg := range pkgs {
		docval := repr.String(*pkg, repr.Indent("\t"))
		fmt.Fprintf(w, "\tcodoc.Register(%s)", docval)
	}
	io.WriteString(w, "}\n")
}
