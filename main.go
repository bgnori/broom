package main

import (
	"./broom"
	"fmt"
	"os"
)

func main() {
	env := broom.NewGlobalRootFrame()
	name := os.Args[1]
	fmt.Printf("%s\n", name)
	var err error
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	if err := broom.Load(file, env, false); err != nil {
		panic(err)
	}
}
