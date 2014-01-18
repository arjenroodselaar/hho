package bytecode

import (
	"testing"
	"strings"
	"io/ioutil"
	"os/exec"
	"os"
)

const testprog1 string = `
package main

func x() {
	return "howdy"
}

func main() {
	println(x())
}
`

const testresult1 string = `.main {
    FPushFuncD 0 "main"
    FCall 0
    PopR # -1, now -1

    Int 0 # 1, now 0
    RetC # -1, now -1
}

.function x() {
    String "howdy" # 1, now 0
    RetC # -1, now -1
    Null # 1, now 0
    RetC # -1, now -1
}

.function main() {
    FPushFuncD 0 "x"
    FCall 0
    UnboxR

    Print
    PopC # -1, now -2
    String "\n" # 1, now -1
    Print
    PopC # -1, now -2
    Null # 1, now -1
    RetC # -1, now -2
}
`

func diffTestFiles(good, bad string) (str string, err error) {
	cleanup := func(f *os.File) {
		f.Close()
		os.Remove(f.Name())
	}
	fgood, err := ioutil.TempFile("/tmp", "hho_test_good")
	defer cleanup(fgood)
	if err != nil {
		return
	}

	fbad, err := ioutil.TempFile("/tmp", "hho_test_bad")
	defer cleanup(fbad)
	if err != nil {
		return
	}

	fgood.WriteString(good)
	fbad.WriteString(bad)

	diffout, _ := exec.Command("diff", "-u", fgood.Name(), fbad.Name()).Output()

	return string(diffout[:]), nil
}

func TestCallInsidePrint(t *testing.T) {
	a := Assemble(testprog1)
	gen := strings.TrimSpace(a.String())
	good := strings.TrimSpace(testresult1)
	diff, _ := diffTestFiles(good, gen)
	if diff != "" {
		t.Error("\n"+diff)
	}
}

func TestSimpleFunction(t *testing.T) {
	code := `package main; func x() int { return 4 }`
	good := `.main {
    FPushFuncD 0 "main"
    FCall 0
    PopR # -1, now -1

    Int 0 # 1, now 0
    RetC # -1, now -1
}

.function x() {
    Int 4 # 1, now 0
    RetC # -1, now -1
}

`
	a := Assemble(code)
	gen := a.String()
	diff, _ := diffTestFiles(good, gen)
	if diff != "" {
		t.Error("\n"+diff)
	}
}
