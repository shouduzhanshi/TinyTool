package build

import (
	"os"
	"tiny_tool/module"
	"tiny_tool/tool"
)

func Build() {
	path := tool.GetCurrentPath()
	appJsonPath := path + "/tiny.json"
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

func buildAndroid(androidDir string) {
	tool.BaseCmd(androidDir+"/gradlew", true, "assembleDebug", "-p", androidDir)
}

func startApp(applicationId string) {
	tool.Adb("shell", "am", "start", "-n", applicationId+"/com.sunmi.android.elephant.core.splash.SplashActivity")
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
