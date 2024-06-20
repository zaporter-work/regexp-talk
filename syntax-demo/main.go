package main

import (
	"fmt"
	"os"
	"regexp/syntax"
)

var flags = syntax.OneLine | syntax.DotNL | syntax.ClassNL | syntax.PerlX | syntax.UnicodeGroups

func main() {

	expr := `a*`
	if len(os.Args) == 2 {
		expr = os.Args[1]
	}
	re, err := syntax.Parse(expr, flags)
	if err != nil {
		panic(err)
	}

	re = re.Simplify()
	prog, err := syntax.Compile(re)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Pattern: %v\n", expr)
	fmt.Println(prog.String())
	fmt.Printf("Is onepass: %v\n", isOnePass(prog))
	fmt.Printf("Cap names: %v\n", re.CapNames())
}

func isOnePass(prog *syntax.Prog) bool {
	return compileOnePass(prog) != nil
}

// compileOnePass returns a new *syntax.Prog suitable for onePass execution if the original Prog
// can be recharacterized as a one-pass regexp program, or syntax.nil if the
// Prog cannot be converted. For a one pass prog, the fundamental condition that must
// be true is: at any InstAlt, there must be no ambiguity about what branch to  take.
func compileOnePass(prog *syntax.Prog) (p *syntax.Prog) {
	if prog.Start == 0 {
		return nil
	}
	// onepass regexp is anchored
	if prog.Inst[prog.Start].Op != syntax.InstEmptyWidth ||
		syntax.EmptyOp(prog.Inst[prog.Start].Arg)&syntax.EmptyBeginText != syntax.EmptyBeginText {
		return nil
	}
	// every instruction leading to InstMatch must be EmptyEndText
	for _, inst := range prog.Inst {
		opOut := prog.Inst[inst.Out].Op
		switch inst.Op {
		default:
			if opOut == syntax.InstMatch {
				return nil
			}
		case syntax.InstAlt, syntax.InstAltMatch:
			if opOut == syntax.InstMatch || prog.Inst[inst.Arg].Op == syntax.InstMatch {
				return nil
			}
		case syntax.InstEmptyWidth:
			if opOut == syntax.InstMatch {
				if syntax.EmptyOp(inst.Arg)&syntax.EmptyEndText == syntax.EmptyEndText {
					continue
				}
				return nil
			}
		}
	}
	return prog
}
