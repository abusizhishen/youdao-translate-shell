package main

import (
	"os"
	"youda-translate-shell/src"
)

func main() {
	src.Login()
}

func query() {
	args := os.Args[1:]
	src.Translate(args[0])
}
