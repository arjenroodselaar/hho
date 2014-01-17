package php

import (
	"fmt"
	"go/ast"
)

type SelectorFunc func([]ast.Expr) []string

func TranslateSelector(target, selector string) SelectorFunc {
	switch target {
	case "fmt":
		return FmtSelector(selector)
	default:
		fmt.Printf("Unsupported target: %s", target)
	}
	return nil
}
