package main

import (
	"./ps"
	"fmt"
)

func main() {
	env := ps.MakeEnv()
	env.Bind("A", ps.Int(42))
	//expr := ps.Int(3)
	//expr := ps.Bool(false)
	//expr := ps.Cons(nil, nil)
	//expr := ps.Cons(ps.Quote, ps.Cons(ps.Int(1), ps.Int(2)))
	//expr := ps.Cons(ps.Name("+"), ps.Cons(ps.Int(1), ps.Int(2)))
	expr := ps.Name("A")
	n := ps.Eval(expr, env)
	fmt.Println(n)
}
