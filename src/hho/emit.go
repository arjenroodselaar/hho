package hho

import (
	"fmt"
	"reflect"
	"go/token"
	"code.google.com/p/go.tools/ssa"
	"strings"
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
	switch t := v.(type) {
	case *ssa.Const:
		switch c := t.Type().(type) {
		case *types.Basic:
			switch c.Kind() {
			case types.String:
				fmt.Printf("\tString %s\n", t.Value)
			default:
				fmt.Println("Unknown Basic type:", c.Kind())
			}
		default:
			fmt.Println("Unknown Const type:", t.Type())
		}
	case *ssa.Builtin:
		switch t.Object().Name() {
		case "print":
			fmt.Printf("\tPrint\n")
			// Pop the 1 pushed on by print.
			fmt.Printf("\tPopC\n")
		default:
			fmt.Printf("Unknown Builtin: %v\n", t.Object().Name())
		}
	default:
		fmt.Println("Unknown Value type:", reflect.TypeOf(t))
	}
}

func EmitCall(i *ssa.Call) {
	for _, value := range(i.Common().Args) {
		EmitValue(value)
	}
	if !i.Common().IsInvoke() {
		EmitValue(i.Common().Value)
	} else {
		fmt.Println("Unknown Call type: invoke")
	}

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
	// Emit the function declaration.
	args := make([]string, len(f.Params))
	for i, param := range(f.Params) {
		args[i] = fmt.Sprintf("$%s", param.Name())
	}
	fmt.Printf(".function %s(%s) {\n", f.String(), strings.Join(args, ", "))
	// Emit the function blocks.
	for _, b := range(f.Blocks) {
		EmitBasicBlock(b)
	}
	// Close.
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
	fmt.Println(".main {")
	// Execute main.init()
	fmt.Printf("\tFPushFuncD 0 \"main.init\"\n\tFCall 0\n\tPopR\n")
	// Execute main.main()
	fmt.Printf("\tFPushFuncD 0 \"main.main\"\n\tFCall 0\n\tPopR\n")
	fmt.Printf("\tNull\n\tRetC\n")
	fmt.Println("}\n")

	for _, pkg := range(p.AllPackages()) {
		EmitPackage(pkg)
	}
}
