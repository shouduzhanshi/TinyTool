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
		build.HotReloadByJavaScript()
	} else if arg == "--build" {
		build.Build()
	} else if arg == "--help" {
		fmt.Println("--hotReload [changeFilePath]")
		fmt.Println("\n热重载")
		fmt.Println("--build [Android Build Project]")
		fmt.Println("\n全量构建")
		fmt.Println("--help")
		fmt.Println("\ntiny_tool命令大全")
		fmt.Println("--create")
		fmt.Println("\n创建工程")
	} else if arg == "--create" {
		project.InitProject()
	} else if arg == "--version" {
		project.InitProject()
	}
}
