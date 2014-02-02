package main

import (
	"go/parser"
	"code.google.com/p/go.tools/importer"
	"code.google.com/p/go.tools/ssa"
	"hho"
	"fmt"
	"bytes"
)
func main() {
	name := "/Users/arjen/dev/hho/examples/for.go"
	imp := importer.New(&importer.Config{})

	// Parse the input file.
	f, err := parser.ParseFile(imp.Fset, name, nil, parser.Mode(0))
	if err != nil {
		panic(err)
	}

	imp.CreatePackage(f.Name.Name, f)
	prog := ssa.NewProgram(imp.Fset, ssa.BuilderMode(0))
	if err = prog.CreatePackages(imp); err != nil {
		panic(err)
	}
	prog.BuildAll()

	//prog.BuildAll()
	//hho.EmitProgram(prog)

	// Create single-file main package and import its dependencies.
	//
	// Create packages for the dependencies.
	//pkg := prog.Package(info.Pkg)
	//pkg.Build()

	//pkg.DumpTo(os.Stdout)

	buf := new(bytes.Buffer)
	translator := hho.NewTranslator(buf)
	translator.EmitProgram(prog)

	//pkg := prog.Package(info.Pkg)
	//pkg.DumpTo(os.Stdout)
	fmt.Println(buf.String())

	//pkg.Func("init").DumpTo(os.Stdout)
	//pkg.Func("main").DumpTo(os.Stdout)
}
