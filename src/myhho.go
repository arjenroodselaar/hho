package main

import (
	"go/parser"
	"go/ast"
	"go/token"
	"fmt"
	"reflect"
	"strings"
)

const hellogo string = `
package main

func printer(x, y int, z string) int {
	print(x + y, z)
	return 1 + 2
}

func main() {
    print("howdy")
    a := 2 + 3
    printer(a)
}
`

func getOpFromKind(t token.Token) string {
	s := ""
	switch t {
	case token.STRING:
		s = "String"
	case token.CHAR:
		s = "String"
	case token.INT:
		s = "Int"
	case token.FLOAT:
		s = "Double"
	case token.ADD:
		s = "Add"
	case token.SUB:
		s = "Sub"
	case token.MUL:
		s = "Mul"
	case token.QUO:
		s = "Div"
	}
	return s
}

type Assembler struct {
	hhas            string
	indent          int
	scopes          []map[string]string
	skip_next_ident bool
}

func NewAssembler() *Assembler {
	a := new(Assembler)
	a.hhas = ""
	a.indent = 0
	a.scopes = make([]map[string]string, 2)
	a.skip_next_ident = false
	return a
}

func (a *Assembler) emit(fstring string, args ...interface{}) (string) {
	ind := strings.Repeat("    ", a.indent)
	str := fmt.Sprintf(fstring, args...)
	return fmt.Sprintf("%s%s\n", ind, str)
}

func (a *Assembler) Print() {
	fmt.Println(a.hhas)
}

func (a *Assembler) EmitArgs(n []ast.Expr) {
	for i, arg := range n {
		switch v := arg.(type) {
		case *ast.BasicLit:
			a.hhas += a.emit("FPassC %d %s", i, v.Value)
		case *ast.Ident:
			a.hhas += a.emit("FPassL %d %s", i, v.Name)
		default:
			fmt.Errorf("Unrecognized type: %s", v)
		}
	}
	a.hhas += a.emit("Fcall %d", len(n))
}

func (a *Assembler) EmitExprStmt(n *ast.ExprStmt) {
	x := n.X;
	switch v := x.(type) {
	case *ast.CallExpr:
		a.hhas += a.emit("FPushFuncD %d %s", len(v.Args), v.Fun.(*ast.Ident).Name)
		a.EmitArgs(v.Args)
	}
}

func (a *Assembler) EmitFuncBody(n *ast.BlockStmt) {
	for _, x := range n.List {
		a.ParseNode(x)
	}
}

func buildArgList(n *ast.FuncDecl) string {
	s := ""
	args := make([]string, 0)
	for _, x := range n.Type.Params.List {
		for _, y := range x.Names {
			args = append(args, y.Name)
		}
	}
	s = strings.Join(args, ", ")
	return s
}

func (a *Assembler) EmitFuncDecl(n *ast.FuncDecl) {
	args := buildArgList(n)
	a.hhas += a.emit(".function %s(%s) {", n.Name.Name, args)
	a.indent++
	a.EmitFuncBody(n.Body)
	a.indent--
	a.hhas += a.emit("}\n\n")
}

func (a *Assembler) EmitBinaryExpr(n *ast.BinaryExpr) {
	//ss := ""
	//TODO: add flag for whether the ident is on the lhs or rhs to generate
	// code accordingly.
}

func (a *Assembler) EmitReturnStmt(n *ast.ReturnStmt) {
	getRetVal := func(e ast.Expr) string {
		switch v := e.(type) {
		case *ast.BinaryExpr:
			x := v.X.(*ast.BasicLit)
			y := v.Y.(*ast.BasicLit)
			s := ""
			s += a.emit(getOpFromKind(x.Kind)+"%s", x.Value)
			s += a.emit(getOpFromKind(v.Y.Kind)+"%s", v.Y.Value)
			s += a.emit(getOpFromKind(v.Op))
			return s
		case *ast.BasicLit:
			return a.emit(getOpFromKind(v.Kind)+"%s", v.Value)
		case *ast.Ident:
			return a.emit("CGetL $%s", v.Name)
		default:
			panic(fmt.Sprintf("Unexpected return: %s\n", reflect.TypeOf(v)))
		}
	}

	if (len(n.Results) == 0) {
		a.hhas += a.emit("RetC")
		return
	}
	if (len(n.Results) == 1) {
		a.hhas += getRetVal(n.Results[0])
	}
	if (len(n.Results) > 1) {
		a.hhas += a.emit("NewArray")
		//		for _, r := range n.Results {

		//		}
	}
	a.hhas += a.emit("RetC")

}

func (a *Assembler) ParseNode(n ast.Node) bool {
	if (n == nil) {
		return false
	}
	switch v := n.(type) {
		//case *ast.Ident:
	case *ast.BinaryExpr:
		a.EmitBinaryExpr(v)
	case *ast.BasicLit:
		a.EmitBasicLit(v)
	case *ast.Ident:
		a.EmitIdent(v)
	case *ast.ReturnStmt:
		a.EmitReturnStmt(v)
	case *ast.FuncDecl:
		a.EmitFuncDecl(v)
		return false
	case *ast.ExprStmt:
		a.EmitExprStmt(v)
	default:
		// v is a ast.Node
		fmt.Println("Not implemented:", reflect.TypeOf(v))
	}
	return true
}

func main() {
	f := token.NewFileSet()
	t, err := parser.ParseFile(f, "hello.go", hellogo, 0)
	if (err != nil) {
		print(err)
	}

	a := NewAssembler()
	ast.Inspect(t, a.ParseNode)
	//ast.Print(f, t)
	a.Print()
}
