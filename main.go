package main

import (
	"fmt"
        "os"
        "./broom"
)

func main() {
	env := broom.NewGlobalRootFrame()
        buf := broom.NewBuffered(os.Stdin)
        expr := broom.BuildSExpr(buf)
        fmt.Println("input:", expr)
        got := broom.Eval(expr, env)
        fmt.Println("-->", got)
}
