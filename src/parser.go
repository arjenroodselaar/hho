package main

import (
  "fmt"
  "reflect"
  "strings"
  "go/parser"
  "go/ast"
  "go/token"
)
func main() {
  fset := token.NewFileSet()

  f, err := parser.ParseFile(fset, "/Users/pwhite/hho/src/examples/main.go", nil, 0)
  if err != nil {
    fmt.Println(err)
    return
  }

  indent := 0

  ast.Inspect(f, func(n ast.Node) bool {
    if n == nil {
      indent--
      fmt.Println(strings.Repeat(" ", indent), "}")
      return false
    }

    v := reflect.ValueOf(n).Elem()
    fmt.Println(strings.Repeat(" ", indent), reflect.TypeOf(n), " {")

    for i := 0; i < v.NumField(); i++ {
      fmt.Println(strings.Repeat(" ", indent), v.Field(i))
    }

    indent++
    return true
  })
}
