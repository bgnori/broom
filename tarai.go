package main

import (
	"fmt"
        "./broom"
        "strings"
)

const tarai = "(define tarai (lambda (x y z) (if (> x y) (tarai (tarai (- x 1) y z) (tarai (- y 1) z x) (tarai (- z 1) x y)) y))) (tarai 3 2 1)"

func main() {
	env := broom.NewGlobalRootFrame()
        buf := broom.NewBuffered(strings.NewReader(tarai))

        for i, expr := range broom.BuildSExpr(buf) {
          fmt.Println("input[", i, "]:", expr)
          got := broom.Eval(expr, env)
          fmt.Println("-->", got)
        }
}

