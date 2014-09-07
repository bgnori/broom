package main

import (
	"fmt"
        "./broom"
        "strings"
)


func main() {
	env := broom.NewGlobalRootFrame()
        buf := broom.NewBuffered(strings.NewReader("(. _env Dump)"))

        for i, expr := range broom.BuildSExpr(buf) {
          fmt.Println("input[", i, "]:", expr)
          got := broom.Eval(expr, env)
          fmt.Println("-->", got)
        }
}

