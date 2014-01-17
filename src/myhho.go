package main

import (
	"bytecode"
	"go/ast"
	"go/parser"
	"go/token"
)

const hellogo string = `
package main

import "fmt"

func main() {
	fmt.Println("hey")
	fmt.Printf("howdy, %s\n", "arjen")
	fmt.Print("sup")
	fmt.Print()
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
