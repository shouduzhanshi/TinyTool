package dev

import (
	"tiny_tool/module"
	"tiny_tool/tool"
)

func Dev() {
	config := tool.GetAppConfig()
	if config.ProjectType == module.JavaScript || config.ProjectType == "" {
		BuildByJavaScript()
	} else if config.ProjectType == module.ES6 || config.ProjectType == module.JSX {
		ByES6()
	}
}
