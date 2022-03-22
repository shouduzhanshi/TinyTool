package main

import (
	"fmt"
	"os"
	"tiny_tool/build"
	"tiny_tool/dev"
	"tiny_tool/log"
	"tiny_tool/project"
)

func main() {
	log.Clean()
	log.Header()
	fmt.Println()
	arg := os.Args[1]
	if arg == "dev" {
		dev.Dev()
	} else if arg == "build" {
		build.Build()
	}  else if arg == "init" {
		project.InitProject()
	} else if arg == "-v" {
		fmt.Println("v1.0.7")
	}else if arg == "-h" {
		fmt.Println("--hotReload [changeFilePath]")
		fmt.Println("\n热重载")
		fmt.Println("--dev [Android Build Project]")
		fmt.Println("\n全量构建")
		fmt.Println("--help")
		fmt.Println("\ntiny_tool命令大全")
		fmt.Println("--create")
		fmt.Println("\n创建工程")
	}
}
