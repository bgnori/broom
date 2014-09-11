package main

import (
  "./broom"
	"fmt"
        "os"
)


func main() {
	env := broom.NewGlobalRootFrame()

        file, err := os.Open("mce.brm")
        if err != nil {
          panic(err)
        }
	buf := broom.NewBuffered(file)
	reader := broom.NewReader(buf)
	builder := broom.NewSExprBuilder()
        for _, expr := range builder.Run(reader).Items() {
          got := broom.Eval(env, expr)
          fmt.Println("-->", got)
        }
        fmt.Println("done.")
/*
	buf = broom.NewBuffered(strings.NewReader(run))
        expr = broom.BuildSExpr(buf)
        got = broom.Eval(env, expr)
        fmt.Println("-->", got)
        */
}
