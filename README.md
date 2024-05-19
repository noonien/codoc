# codoc
[![GoDoc](https://pkg.go.dev/badge/github.com/noonien/codoc)](https://pkg.go.dev/github.com/noonien/codoc)

Extract Go package documentation to be used at runtime.


# Usage
To generate a documentation file for a Go package, run the following command:

```shell
go run github.com/noonien/codoc/cmd/codoc@latest -pkg main -out example_doc.go ./example
```

This will generate an `example_doc.go` file that registers the package documentation with the `codoc` package.
Assuming the `example_doc.go` file is imported somewhere by your program, the documentation can be accessed like so:

```go
package main

import (
    "github.com/noonien/codoc"

    // make sure example_doc.go is part of, or imported by your program
)

func main() {
    pkg := codoc.Package("main")      // Get documentation for the main package
    fn := codoc.Function("main.Foo")  // Get documentation for the function Foo in the main package
    st := codoc.Struct("main.Bar")    // Get documentation for the struct Bar in the main package

    // Example usage
    fmt.Println(pkg)
    fmt.Println(fn)
    fmt.Println(st)
}
```


# License
`codoc` is released under the MIT License. See the `LICENSE` file for more details.