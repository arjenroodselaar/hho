package hho

import (
	"io"
	"fmt"
	"strings"
	"code.google.com/p/go.tools/ssa"
	"code.google.com/p/go.tools/go/types"
	"go/token"
)

type Translator struct {
	writer io.Writer
}

func NewTranslator(w io.Writer) *Translator {
	return &Translator{w}
}

func (t *Translator) EmitLabel(label string) {
	fmt.Fprintf(t.writer, "%s:\n", label)
}

func (t *Translator) EmitStatement(stmt string) {
	fmt.Fprintf(t.writer, "\t%s\n", stmt)
}

func (t *Translator) EmitReturn(i *ssa.Return) {
	if len(i.Results) == 0 {
		// Push Null on the stack
		t.EmitStatement("Null")
	}
	t.EmitStatement("RetC")
}

func (t *Translator) EmitIf(i *ssa.If) {
	t.EmitValue(i.Cond)
	// Skip the next basic block if not True, continue otherwise.
	t.EmitStatement(fmt.Sprintf("JmpNZ %s.%s", i.Parent().String(), i.Block().Succs[0]))
}

func (t *Translator) EmitJump(i *ssa.Jump) {
	t.EmitStatement(fmt.Sprintf("Jmp %s.%s", i.Parent().String(), i.Block().Succs[0]))
}

func (t *Translator) EmitRegisterLoad(v ssa.Value) {
	t.EmitStatement(fmt.Sprintf("CGetL $%s", v.Name()))
}

func (t *Translator) EmitUnOp(i *ssa.UnOp) {
	t.EmitValue(i.X)

	switch i.Op {
	case token.MUL:
		// Pointer indirection, load either a local or a global here
		switch i.X.(type) {
		case *ssa.Global:
			t.EmitStatement("CGetG")
		default:
			t.EmitStatement("CGetN")
		}
	default:
		panic(fmt.Errorf("Unknown UnOp (%s)", i.Op))
	}

	t.EmitStatement(fmt.Sprintf("SetL $%s", i.Name()))
	t.EmitStatement("PopC")
}

func (t *Translator) EmitBinOp(i *ssa.BinOp) {
	// Load the parameters and push them onto the stack.
	t.EmitValue(i.X)
	t.EmitValue(i.Y)

	switch i.Op {
	case token.ADD:
		t.EmitStatement("Add")
	case token.MUL:
		t.EmitStatement("Mul")
	case token.QUO:
		t.EmitStatement("Div")
	case token.LSS:
		t.EmitStatement("Le")
	default:
		panic(fmt.Errorf("Unknown BinOp (%s)", i.Op))
	}
	t.EmitStatement(fmt.Sprintf("SetL $%s", i.Name()))
	t.EmitStatement("PopC")
}

func (t *Translator) EmitStore(i *ssa.Store) {
	t.EmitValue(i.Addr)
	t.EmitValue(i.Val)
	switch i.Addr.(type) {
	case *ssa.Global:
		t.EmitStatement("SetG")
	default:
		t.EmitStatement("SetN")
	}
	t.EmitStatement("PopC")
}

func (t *Translator) EmitFunctionCall(i *ssa.Call) {
	switch typed := i.Common().Value.(type) {
	case *ssa.Function:
		t.EmitStatement(fmt.Sprintf("FPushFuncD %d \"%s\"", len(typed.Params), typed.String()))
		for i, arg := range(i.Common().Args) {
			t.EmitValue(arg)
			t.EmitStatement(fmt.Sprintf("FPassC %d", i))
		}
		t.EmitStatement(fmt.Sprintf("FCall %d", len(i.Common().Args)))
		t.EmitStatement("UnboxR")
	case *ssa.Builtin:
		for _, arg := range(i.Common().Args) {
			t.EmitValue(arg)
		}
		t.EmitValue(typed)
	default:
		panic(fmt.Errorf("Unknown Call (type: %T)", typed))
	}
	t.EmitStatement(fmt.Sprintf("SetL $%s", i.Name()))
	t.EmitStatement("PopC")
}

func (t *Translator) EmitValue(v ssa.Value) {
	switch typedValue := v.(type) {
	case *ssa.Const:
		switch constValue := v.Type().(type) {
		case *types.Basic:
			switch constValue.Kind() {
			case types.Bool:
				t.EmitStatement(strings.Title(typedValue.Value.String()));
			case types.Int:		fallthrough
			case types.Int8:	fallthrough
			case types.Int16:	fallthrough
			case types.Int32:	fallthrough
			case types.Int64:
				t.EmitStatement(fmt.Sprintf("Int %s", typedValue.Value));
			case types.String:
				t.EmitStatement(fmt.Sprintf("String %s", typedValue.Value));
			default:
				panic(fmt.Errorf("Unknown Basic (%s)", constValue.Kind()))
			}
		default:
			panic(fmt.Errorf("Unknown Const (type %T)", constValue))
		}
	case *ssa.Builtin:
		switch typedValue.Object().Name() {
		case "print":
			t.EmitStatement("Print")
		case "println":
			t.EmitStatement("Print")
			t.EmitStatement("PopC")
			t.EmitStatement("String \"\\n\"")
			t.EmitStatement("Print")
		default:
			panic(fmt.Errorf("Unknown Builtin (%s)", typedValue.Object().Name()))
		}
	case *ssa.Global:
		t.EmitStatement(fmt.Sprintf("String \"%s\"", typedValue.Name()))
	case *ssa.Call:
		t.EmitRegisterLoad(v)
	case *ssa.Parameter:
		t.EmitRegisterLoad(v)
	case *ssa.UnOp:
		t.EmitRegisterLoad(v)
	case *ssa.BinOp:
		t.EmitRegisterLoad(v)
	default:
		panic(fmt.Errorf("Unknown Value (type %T, %s)", typedValue, typedValue))
	}
}

func (t *Translator) EmitInstruction(i ssa.Instruction) {
	t.EmitStatement(fmt.Sprintf("# %s", i.String()))
	switch typedInstruction := i.(type) {
	// Control flow instructions
	case *ssa.If:
		t.EmitIf(typedInstruction);
	case *ssa.Jump:
		t.EmitJump(typedInstruction)
	case *ssa.Return:
		t.EmitReturn(typedInstruction)
	// Operators
	case *ssa.UnOp:
		t.EmitUnOp(typedInstruction)
	case *ssa.BinOp:
		t.EmitBinOp(typedInstruction)
	case *ssa.Store:
		t.EmitStore(typedInstruction)
	// Function instructions
	case *ssa.Call:
		t.EmitFunctionCall(typedInstruction)
	default:
		panic(fmt.Errorf("Unknown Instruction (type %T, %s)", i, i))
	}
}

func (t *Translator) EmitBasicBlock(b *ssa.BasicBlock) {
	t.EmitLabel(fmt.Sprintf("%s.%s", b.Parent(), b.String()))
	for _, instr := range(b.Instrs) {
		t.EmitInstruction(instr);
	}
}

func (t *Translator) EmitFunctionDefinition(f *ssa.Function) {
	// Emit the function declaration.
	args := make([]string, len(f.Params))
	for i, param := range(f.Params) {
		args[i] = fmt.Sprintf("$%s", param.Name())
	}
	fmt.Fprintf(t.writer, ".function %s(%s) {\n", f.String(), strings.Join(args, ", "))
	// Emit the function blocks.
	for _, b := range(f.Blocks) {
		t.EmitBasicBlock(b)
	}
	// Close.
	fmt.Fprintln(t.writer, "}\n");
}

func (t *Translator) EmitPackage(pkg *ssa.Package) {
	for _, member := range(pkg.Members) {
		switch m := member.(type) {
		case *ssa.Function:
			t.EmitFunctionDefinition(m)
		case *ssa.Global:
			// Globals are taken care of in the init function of packages, so ignore here.
		default:
			panic(fmt.Errorf("Unknown Package (type %T)", m))
		}
	}
}

func (t *Translator) EmitProgram(program *ssa.Program) {
	// Emit the .main directive
	fmt.Fprintln(t.writer, ".main {")
	// Init global main.init$guard
	t.EmitStatement("String \"main.init$guard\"")
	t.EmitStatement("False")
	t.EmitStatement("SetG")
	t.EmitStatement("PopC")
	// Call main.init()
	t.EmitStatement("FPushFuncD 0 \"main.init\"")
	t.EmitStatement("FCall 0")
	t.EmitStatement("PopR")
	// Call main.main()
	t.EmitStatement("FPushFuncD 0 \"main.main\"")
	t.EmitStatement("FCall 0")
	t.EmitStatement("PopR")
	// Return
	t.EmitStatement("Null")
	t.EmitStatement("RetC")
	fmt.Fprintln(t.writer, "}\n")

	for _, pkg := range(program.AllPackages()) {
		t.EmitPackage(pkg)
	}
}
