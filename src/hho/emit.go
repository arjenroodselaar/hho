package hho

import (
	"fmt"
	"reflect"
	"go/token"
	"code.google.com/p/go.tools/ssa"
	"code.google.com/p/go.tools/go/types"
)

func EmitIf(i *ssa.If) {
	//fmt.Printf("\tJmpNZ %s (%s)\n", i.Block().Succs[0], i)
}

func EmitUnOp(o *ssa.UnOp) {
	switch o.Op {
	case token.MUL:
		// Ignore
	default:
		fmt.Println("Unknown UnOp:", o.Op)
	}
}

func EmitJump(j *ssa.Jump) {
	fmt.Printf("\tJmp %s\n", j.Block().Succs[0])
}

func EmitReturn(r *ssa.Return) {
	if len(r.Results) == 0 {
		// Push Null on the stack
		fmt.Printf("\tNull\n")
	}
	fmt.Printf("\tRetC\n")
}

func EmitValue(v ssa.Value) {
	switch t := v.Type().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.String:
			fmt.Printf("\tString %s\n", "bla")
		}
	//case *types.Tuple:
	//	fmt.Println("tuple len", v.(*types.Tuple).Len())
	default:
		fmt.Println("Unknown Value type:", reflect.TypeOf(t))
	}
}

func EmitCall(i *ssa.Call) {
	for _, value := range(i.Call.Args) {
		EmitValue(value)
	}
	fmt.Println(i.Common().Description())
	EmitValue(i.Common().Value)
}

func EmitInstruction(i ssa.Instruction) {
	switch t := i.(type) {
	case *ssa.If:
		EmitIf(t)
	case *ssa.UnOp:
		EmitUnOp(t)
	case *ssa.Store:
		// Ignore
	case *ssa.Jump:
		EmitJump(t)
	case *ssa.Return:
		EmitReturn(t)
	case *ssa.Call:
		EmitCall(t)
	default:
		fmt.Println("Unknown Instruction:", reflect.TypeOf(t))
	}
}

func EmitBasicBlock(b *ssa.BasicBlock) {
	fmt.Printf("%s:\n", b.String())
	for _, instr := range(b.Instrs) {
		EmitInstruction(instr)
	}
}

func EmitFunction(f *ssa.Function) {
	fmt.Printf(".function %s() {\n", f.String())
	for _, b := range(f.Blocks) {
		EmitBasicBlock(b)
	}
	fmt.Printf("}\n\n")
}

func EmitGlobal(g *ssa.Global) {}

func EmitPackage(p *ssa.Package) {
	for _, member := range(p.Members) {
		switch t := member.(type) {
		case *ssa.Function:
			EmitFunction(t)
		case *ssa.Global:
			EmitGlobal(t)
		default:
			fmt.Println("Unknown Member type: ", reflect.TypeOf(t));
		}
	}
}

func EmitProgram (p *ssa.Program) {
	for _, pkg := range(p.AllPackages()) {
		EmitPackage(pkg)
	}
}
