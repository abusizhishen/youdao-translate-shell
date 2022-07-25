package main

import (
	"fmt"
	"os"
	"youda-translate-shell/src"
)

var args []string

func init() {
	args = os.Args[1:]
}

func main() {

	if len(args) == 0 {
		fmt.Println("请输入查询单词")
		return
	}
	src.Login(args[0])
	return
	query()
}

func query() {

	src.Translate(args[0])
	info, err := src.CheckLoginStatus()
	fmt.Println(err, info)
}
