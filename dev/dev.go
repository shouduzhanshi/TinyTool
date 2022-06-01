package dev

import (
	"os"
	"os/exec"
	"tiny_tool/module"
	"tiny_tool/tool"
)

var projectPath = tool.GetCurrentPath()

var cmd *exec.Cmd

func Dev() {
	config := tool.GetAppConfig()
	if config.ProjectType == module.JavaScript || config.ProjectType == "" {
		BuildByJavaScript()
	} else if config.ProjectType == module.ES6 || config.ProjectType == module.JSX {
		ByES6()
	}
}

func WebpackDev(initCallback func()) {
	mute := false
	if os.Getenv("SHOW_BUILD_LOG") == "false"{
		mute = true
	}
	cmd = tool.CmdWatch(mute,initCallback, "npm", "run", "watch", "--prefix", projectPath)
}
