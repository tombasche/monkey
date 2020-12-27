package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/compiler/vm"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	macroEnv := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		comp := compiler.New()
		err := comp.Compile(expanded)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s \n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed: \n %s \n", err)
			continue
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Parser errors: \n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
