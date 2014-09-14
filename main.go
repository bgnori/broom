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
	test := strings.Contains(os.Args[1], "t")
	for i, name := range os.Args[2:] {
		if file, err := os.Open(name); err != nil {
			panic(err)
		} else {
			v := verbose || (test && i == len(os.Args)-3)
			if err := broom.Load(file, env, v); err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", name)
		}
	}
	broom.Repl(os.Stdin, env) //, os.Stdout)
}
