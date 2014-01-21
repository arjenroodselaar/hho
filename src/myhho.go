package main

import (
	"bytecode"
	"go/ast"
	"go/parser"
	"go/token"
)

const hellogo string = `
package main

func main() {
	a := make([]string, 2)
	append(a, "howdy")
	append(a, "bass")
	append(a, "man")
	var_dump(a)
}
`

func main() {
	f := token.NewFileSet()
	t, err := parser.ParseFile(f, "hello.go", hellogo, 0)

	if (err != nil) {
		panic(err)
	}


	a := bytecode.NewAssembler()
	ast.Inspect(t, a.ParseNode)
	ast.Print(f, t)
	a.Print()
}
