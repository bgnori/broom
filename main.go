package main

import (
	"./broom"
	"os"
)

func main() {
	broom.Repl(os.Stdin) //, os.Stdout)
}
