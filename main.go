package main

import (
	"./broom"
	"fmt"
	"os"
	"strings"
)

func main() {
	env := broom.NewGlobalRootFrame()

	verbose := strings.Contains(os.Args[1], "V")

	for _, name := range os.Args[2:] {
		if file, err := os.Open(name); err != nil {
			panic(err)
		} else {
			if err := broom.Load(file, env, verbose); err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", name)
		}
	}
	broom.Repl(os.Stdin, env) //, os.Stdout)
}
