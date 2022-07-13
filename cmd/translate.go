package main

import (
	"os"
	"youda-translate-shell/src"
)

func main() {
	args := os.Args[1:]
	src.Translate(args[0])
}
