package main

import (
	"./broom"
	"fmt"
	"os"
)

func main() {
	env := broom.NewGlobalRootFrame()
	for _, name := range os.Args[1:] {
		if file, err := os.Open(name); err != nil {
			panic(err)
		} else {
			if err := broom.Load(file, env); err != nil {
				panic(err)
			}
			fmt.Printf("loaded %s\n", name)
		}
	}
	broom.Repl(os.Stdin, env) //, os.Stdout)
}
