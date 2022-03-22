package main

import (
	"fmt"
	"os"
	"tiny_tool/build"
	"tiny_tool/project"
)

func main() {
	arg := os.Args[1]
	if arg == "--hotReload" {
		build.HotReloadByJavaScript()
	} else if arg == "dev" {
		build.Build()
	} else if arg == "-h" {
		fmt.Println("--hotReload [changeFilePath]")
		fmt.Println("\n热重载")
		fmt.Println("--build [Android Build Project]")
		fmt.Println("\n全量构建")
		fmt.Println("--help")
		fmt.Println("\ntiny_tool命令大全")
		fmt.Println("--create")
		fmt.Println("\n创建工程")
	} else if arg == "init" {
		project.InitProject()
	} else if arg == "-v" {
		fmt.Println("v1.0.7")
	}
}
