package main

import (
	"MockConfig/build"
	"MockConfig/install"
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
	} else if arg == "--install" {
		install.Install()
	} else if arg == "--help" {
		fmt.Println("--hotReload")
		fmt.Println("热重载")
		fmt.Println("--build")
		fmt.Println("全量构建")
		fmt.Println("--install")
		fmt.Println("安装tiny_tool")
		fmt.Println("--help")
		fmt.Println("tiny_tool命令大全")
		fmt.Println("--create")
		fmt.Println("create")
	} else if arg == "--create" {
		project.InitProject()
	}else if arg == "--v" {
		fmt.Println("version:0.0.1")
		installPath:=os.Getenv("TINY_TOOL")+"/tiny_tool"
		fmt.Println("install path:",installPath)
	}
}
