package build

import (
	"os"
	"tiny_tool/module"
	"tiny_tool/tool"
)

func Build() {
	path := tool.GetCurrentPath()
	appJsonPath := path + "/"+module.TINY_JSON
	json := tool.DeCodeAppJson(appJsonPath)
	if json.ProjectType == module.JavaScript || json.ProjectType == "" {
		BuildByJavaScript(false)
	} else if json.ProjectType == module.ES6 {
		ByES6(path, json)
	}
}

func installAndroidApp(androidDir string) (int, error) {
	return tool.Adb("install", "-r", androidDir+"/build/outputs/apk/debug/app-debug.apk")
}

func buildAndroid(androidDir string) int {
	 cmd, _, _ := tool.BaseCmd(androidDir+"/gradlew", false, "assembleDebug", "-p", androidDir)
	return cmd
}

func startApp(applicationId string) {
	tool.Adb("shell", "am", "start", "-n", applicationId+"/com.whl.tinyui.sample.MainActivity")
}

func getElephantDir(isHotReload bool) string {
	if isHotReload {
		if len(os.Args) >= 4 {
			return os.Args[3]
		}
	} else {
		if len(os.Args) >= 3 {
			return os.Args[2]
		}
	}
	return ""
}
