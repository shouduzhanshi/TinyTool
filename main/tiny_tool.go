package main

import (
	"MockConfig/build"
	"MockConfig/project"
	"fmt"
	"os"
)

func main() {
	arg := os.Args[1]
	if arg == "--hotReload" {
		build.HotReload()
	} else if arg == "--build" {
		build.Build()
	} else if arg == "--help" {
		fmt.Println("--hotReload [changeFilePath]")
		fmt.Println("热重载")
		fmt.Println("--build")
		fmt.Println("全量构建")
		fmt.Println("--help")
		fmt.Println("tiny_tool命令大全")
		fmt.Println("--create")
		fmt.Println("create")
	} else if arg == "--create" {
		project.InitProject()
	}
}
