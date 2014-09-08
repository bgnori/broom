package broom

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Repl(in io.Reader) {
	env := NewGlobalRootFrame()
	env.Bind("dump", Closure(func(dynamic Enviroment, cdr Pair) Value {
		dynamic.Dump()
		return nil
	}))

	reader := bufio.NewReader(in)

	chunks := make([]string, 0)

	fmt.Println("Hello!")
	fmt.Print("broom > ")

	for line, _, err := reader.ReadLine(); err != io.EOF; line, _, err = reader.ReadLine() {
		if err != nil {
			panic(err)
		}
		if len(line) == 0 {
			c := strings.Join(chunks, "")
			chunks = chunks[0:0]
			expr := BuildSExpr(NewBuffered(strings.NewReader(c)))
			fmt.Println("input:", expr)
			got := Eval(expr, env)
			fmt.Println("-->", got)
			fmt.Print("broom > ")
		} else {
			fmt.Println("... ")
			chunks = append(chunks, string(line))
		}
	}
	fmt.Println()
	fmt.Println("bye!")
}
