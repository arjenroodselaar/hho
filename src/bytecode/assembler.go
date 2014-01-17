package bytecode

import (
	"fmt"
	"go/ast"
	"go/token"
	"php"
	"reflect"
	"strings"
)

type Assembler struct {
	hhas            string

	cur_label,
	indent,
	stack_count    int

	in_assign,
	in_lhs,
	skip_next_ident,
	trace_stack        bool
}

func NewAssembler() *Assembler {
	a := new(Assembler)
	a.hhas = ""
	a.indent = 0
	a.cur_label = 0
	a.stack_count = 0
	a.skip_next_ident = false
	a.trace_stack = true
	return a
}

func (a *Assembler) emit(fstring string, args ...interface{}) (string) {
	b := ""
	if (a.trace_stack) {
		op := strings.Fields(fstring)[0]
		x := LookupStackDelta(op)
		if (x != 0) {
			a.stack_count += x
			b = fmt.Sprintf(" # %d, now %d", x, a.stack_count)
		}
	}
	ind := strings.Repeat("    ", a.indent)
	str := fmt.Sprintf(fstring, args...)
	return ind + str + b + "\n"
}

func (a *Assembler) emitLabel(l string) {
	a.hhas += a.emit("%s:", l)
}

func (a *Assembler) emitMultiple(s []string) {
	for _, n := range s {
		a.hhas += a.emit(n)
	}
}

func (a *Assembler) Print() {
	fmt.Println(a.hhas)
}

func buildArgList(n *ast.FuncDecl) string {
	s := ""
	args := make([]string, 0)
	for _, x := range n.Type.Params.List {
		for _, y := range x.Names {
			args = append(args, "$"+y.Name)
		}
	}
	s = strings.Join(args, ", ")
	return s
}

func (a *Assembler) EmitAssignStmt(n *ast.AssignStmt) {
	a.in_assign = true
	a.ParseNode(n.Rhs[0])
	a.in_lhs = true
	a.ParseNode(n.Lhs[0])
	a.in_lhs = false
	a.in_assign = false
}

func (a *Assembler) EmitBasicLit(n *ast.BasicLit) {
	a.hhas += a.emit(LookupOpFromKind(n.Kind)+" %s", n.Value)
}

func (a *Assembler) EmitBinaryExpr(n *ast.BinaryExpr) {
	a.ParseNode(n.Y)
	a.ParseNode(n.X)
	a.hhas += a.emit(LookupOpFromKind(n.Op))

}

func (a *Assembler) EmitBlockStmt(n *ast.BlockStmt) {
	for _, x := range n.List {
		a.ParseNode(x)
	}
}

func (a *Assembler) EmitCallArgs(n []ast.Expr) {
	for i, arg := range n {
		switch v := arg.(type) {
		case *ast.BasicLit:
			a.hhas += a.emit("%s %s", LookupOpFromKind(v.Kind), v.Value)
			a.hhas += a.emit("FPassC %d", i)
		case *ast.Ident:
			a.hhas += a.emit("FPassL %d $%s", i, v.Name)
		default:
			fmt.Printf("Unrecognized type: %s\n", v)
		}
	}
	a.hhas += a.emit("FCall %d", len(n))
	a.hhas += a.emit("PopC")
}

func (a *Assembler) EmitCallExpr(n *ast.CallExpr) {
	fname := ""
	emitter := php.SelectorFunc(nil)
	switch n.Fun.(type) {
	case *ast.SelectorExpr:
		sel := n.Fun.(*ast.SelectorExpr)
		emitter = php.TranslateSelector(sel.X.(*ast.Ident).Name, sel.Sel.Name)
	case *ast.Ident:
		fname = n.Fun.(*ast.Ident).Name
		emitter = func(args []ast.Expr) []string {
			return []string{fmt.Sprintf("FPushFuncD %d \"%s\"", len(args), fname)}
		}
	}

	lines := emitter(n.Args)
	a.emitMultiple(lines)
	a.EmitCallArgs(n.Args)
}

func (a *Assembler) EmitFile(n *ast.File) {
	//TODO: Bunch of stuff relating to packages, etc
	main := new(ast.CallExpr)
	main.Fun = new(ast.Ident)
	main.Fun.(*ast.Ident).Name = "main"
	a.hhas += a.emit(".main {")
	a.indent++
	a.EmitCallExpr(main)
	a.hhas += a.emit("RetC")
	a.indent--
	a.hhas += a.emit("}\n")
	a.skip_next_ident = true
}

func (a *Assembler) EmitForStmt(n *ast.ForStmt) {
	a.ParseNode(n.Init)
	label := a.getNextLabel()
	a.emitLabel(label + "_for")
	a.indent++

	a.ParseNode(n.Cond)
	a.hhas += a.emit("JmpNZ %s", label+"_end")

	a.indent--
	a.emitLabel(label + "_loop")
	a.indent++

	a.ParseNode(n.Body)

	a.indent--
	a.emitLabel(label + "_post")
	a.indent++

	a.ParseNode(n.Post)
	a.hhas += a.emit("Jmp %s", label+"_for")

	a.indent--
	a.emitLabel(label + "_end")
}

func (a *Assembler) EmitFuncBody(n *ast.BlockStmt) {
	for _, x := range n.List {
		a.ParseNode(x)
	}
}

func (a *Assembler) EmitFuncDecl(n *ast.FuncDecl) {
	args := buildArgList(n)
	a.hhas += a.emit(".function %s(%s) {", n.Name.Name, args)
	a.indent++
	a.EmitFuncBody(n.Body)
	if (n.Type.Results == nil) {
		a.hhas += a.emit("Null")
		a.hhas += a.emit("RetC")
	}
	a.indent--
	a.hhas += a.emit("}\n")
}

func (a *Assembler) EmitIdent(n *ast.Ident) {
	if (a.in_assign) {
		if (a.in_lhs) {
			a.hhas += a.emit("SetL $%s", n.Name)
			a.hhas += a.emit("PopC")
			return
		}
	}
	a.hhas += a.emit("CGetL $%s", n.Name)
}

func (a *Assembler) EmitIfStmt(n *ast.IfStmt) {
	a.ParseNode(n.Cond)
	label := a.getNextLabel()
	elseLabel := label + "_else"
	endLabel := label + "_end"

	a.hhas += a.emit("JmpNZ %s", elseLabel)
	a.emitLabel(label)
	a.indent++
	a.ParseNode(n.Body)
	a.hhas += a.emit("Jmp %s", endLabel)
	a.indent--
	a.emitLabel(elseLabel)
	a.indent++
	a.ParseNode(n.Else)
	a.indent--
	a.emitLabel(endLabel)
}

func (a *Assembler) EmitIncDecStmt(n *ast.IncDecStmt) {
	op := ""
	switch n.Tok {
	case token.INC:
		op = "PostInc"
	case token.DEC:
		op = "PostDec"
	default:
		panic(fmt.Sprintf("Whoa shit, unrecognized IncDec %d", n.Tok))
	}

	a.hhas += a.emit("IncDecL $%s %s", n.X.(*ast.Ident).Name, op)
	a.hhas += a.emit("PopC")
}

func (a *Assembler) EmitParenExpr(n *ast.ParenExpr) {
	a.ParseNode(n.X)
}

func (a *Assembler) getNextLabel() (lbl string) {
	lbl = fmt.Sprintf("label_%d", a.cur_label)
	a.cur_label++
	return
}

func (a *Assembler) EmitReturnStmt(n *ast.ReturnStmt) {
	getRetVal := func(e ast.Expr) string {
		switch v := e.(type) {
		case *ast.BinaryExpr:
			x := v.X.(*ast.BasicLit)
			y := v.Y.(*ast.BasicLit)
			s := ""
			s += a.emit(LookupOpFromKind(x.Kind)+"%s", x.Value)
			s += a.emit(LookupOpFromKind(y.Kind)+"%s", y.Value)
			s += a.emit(LookupOpFromKind(v.Op))
			return s
		case *ast.BasicLit:
			return a.emit(LookupOpFromKind(v.Kind)+"%s", v.Value)
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

func (a *Assembler) EmitSelectorExpr(n *ast.SelectorExpr) {
	//fmt.Printf("%#v\n", n)
}

func (a *Assembler) ParseNode(n ast.Node) bool {
	if (n == nil) {
		return false
	}
	switch v := n.(type) {
	case *ast.AssignStmt:
		a.EmitAssignStmt(v)
	case *ast.BinaryExpr:
		a.EmitBinaryExpr(v)
		return false
	case *ast.BasicLit:
		a.EmitBasicLit(v)
	case *ast.BlockStmt:
		a.EmitBlockStmt(v)
	case *ast.CallExpr:
		a.EmitCallExpr(v)
		return false
	case *ast.ExprStmt:
		a.ParseNode(v.X)
	case *ast.File:
		a.EmitFile(v)
	case *ast.ForStmt:
		a.EmitForStmt(v)
	case *ast.FuncDecl:
		a.EmitFuncDecl(v)
		return false
	case *ast.Ident:
		if (!a.skip_next_ident) {
			a.EmitIdent(v)
		} else {
			a.skip_next_ident = false
		}
	case *ast.IfStmt:
		a.EmitIfStmt(v)
	case *ast.IncDecStmt:
		a.EmitIncDecStmt(v)
	case *ast.ImportSpec:
		return false
	case *ast.ParenExpr:
		a.EmitParenExpr(v)
	case *ast.ReturnStmt:
		a.EmitReturnStmt(v)
		return false
	case *ast.SelectorExpr:
		a.EmitSelectorExpr(v)
	default:
		// v is a ast.Node
		fmt.Println("Not implemented:", reflect.TypeOf(v))
	}
	return true
}
