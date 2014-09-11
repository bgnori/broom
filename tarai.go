package main

import (
	"./broom"
	"fmt"
	"strings"
)

const tarai = "(define tarai (lambda (x y z) (if (> x y) (tarai (tarai (- x 1) y z) (tarai (- y 1) z x) (tarai (- z 1) x y)) y)))"
const run = "(tarai 3 2 1)"

func main() {
	env := broom.NewGlobalRootFrame()

	buf := broom.NewBuffered(strings.NewReader(tarai))
	expr := broom.BuildSExpr(buf)
	got := broom.Eval(env, expr)
	fmt.Println("-->", got)

	buf = broom.NewBuffered(strings.NewReader(run))
	expr = broom.BuildSExpr(buf)
	got = broom.Eval(env, expr)
	fmt.Println("-->", got)
}
