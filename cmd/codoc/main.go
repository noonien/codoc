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

var (
	outFile  = flag.String("out", "", "output file, leave empty to write to stdout")
	pkgName  = flag.String("pkg", "", "output file package")
	exported = flag.Bool("e", false, "only register exported functions and structs")
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	flag.Parse()
	if len(*pkgName) == 0 {
		flag.Usage()
		log.Fatal("missing flag: pkg")
	}

	paths := flag.Args()
	if len(paths) == 0 {
		flag.Usage()
		log.Fatalf("no package paths specified")
	}

	opts := []codocgen.Option{}
	if *exported {
		opts = append(opts, codocgen.Exported())
	}

	var pkgs []*codoc.Package
	for _, p := range flag.Args() {
		pkg, err := codocgen.FromPath(p, opts...)
		if err != nil {
			log.Fatalf("could not get docs for %q: %v", p, err)
		}
		log.Printf("got docs for %s", pkg.Name)
		pkgs = append(pkgs, pkg)
	}

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

func writeDoc(w io.WriteCloser, pkgs []*codoc.Package) {
	defer w.Close()

	fmt.Fprintf(w, "// generated @ %s by gendoc\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(w, "package %s\n", *pkgName)
	fmt.Fprintln(w)
	io.WriteString(w, "import \"github.com/noonien/codoc\"\n")
	fmt.Fprintln(w)

	io.WriteString(w, "func init() {\n")
	for _, pkg := range pkgs {
		docval := repr.String(*pkg, repr.Indent("\t"))
		fmt.Fprintf(w, "\tcodoc.Register(%s)", docval)
	}
	io.WriteString(w, "}\n")
}
