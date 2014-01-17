package php

import (
	"fmt"
	"go/ast"
)

/*
func printargs(args []ast.Expr) {
	for _, arg := range args {
		ast.Inspect(arg, a.ParseNode)
		a.hhas += a.emit("Print")
		a.hhas += a.emit("PopC")
	}
}

if (fname == "print" || fname == "println") {
	printargs(n.Args)
	if (fname == "println") {
		a.hhas += a.emit("String \"\\n\"")
		a.hhas += a.emit("Print")
		a.hhas += a.emit("PopC")
	}
}*/

func Print(args []ast.Expr) []string {
	return []string{fmt.Sprintf("FPushFuncD %d \"print\"", len(args))}
}

func Println(args []ast.Expr) []string {
	return []string{fmt.Sprintf("FPushFuncD %d \"println\"", len(args))}
}

func Printf(args []ast.Expr) []string {
	fmt.Printf("%#v", args)
	return []string{fmt.Sprintf("FPushFuncD %d \"printf\"", len(args))}
}

func FmtSelector(selector string) SelectorFunc {
	switch selector {
	case "Println":
		return Println
	case "Print":
		return Print
	case "Printf":
		return Printf
	}
	return nil
}
