package main

import (
	"fmt"
        "os"
        "./ps"
)

func main() {
	env := ps.MakeEnv()
	env.Bind("+", ps.BuiltinPlus)
        for _, expr := range ps.Parse(os.Stdin) {
            got := ps.Eval(expr, env)
            fmt.Println("-->", got)
        }
}
