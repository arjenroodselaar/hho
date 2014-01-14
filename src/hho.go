package main

import (
	"go/token"
	"go/parser"
	"go/ast"
)

func get_ast(name string) (fset *token.FileSet, f *ast.File, err error) {
	fset = token.NewFileSet()
	f, err = parser.ParseFile(fset, name, nil, parser.Mode(0))
	return
}

func print_ast(fset *token.FileSet, f *ast.File) {
	ast.Print(fset, f)
}

func print_hhas(fset *token.FileSet, f *ast.File) {
}

func main() {
	fset, f, err := get_ast("/Users/arjen/dev/hho/examples/hello_world.go")
	if (err != nil) {
		panic(err)
	}
	print_ast(fset, f)
	print_hhas(fset, f)
}
