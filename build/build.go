package build

import (
	"MockConfig/tool"
	"fmt"
	"os"
)

func HotReload() {
	path := tool.GetCurrentPath()
	appJson := path + "/tiny.json"
	json := tool.DeCodeAppJson(appJson)

	applicationId := tool.GetApplicationId(json)
	if adb, _ := tool.Adb("shell", "pm", "path", applicationId); adb == 1 {
		Build()
	}
	changeFilePath := os.Args[2]
	tool.Adb("shell", "mkdir", "sdcard/Android/data/"+applicationId)
	tool.Adb("shell", "mkdir", "sdcard/Android/data/"+applicationId+"/cache")
	tool.Adb("push", changeFilePath, "sdcard/Android/data/"+applicationId+"/cache")
	if file, err := os.Open(changeFilePath); err == nil {
		defer file.Close()
		if info, err := file.Stat(); err == nil {
			tool.Adb("shell", "am", "start", "-n", applicationId+"/com.sunmi.android.elephant.api.container.ContainerActivity", "--es", "hotReLoad", "\"hotReLoad://cache/"+info.Name()+"\"")
		}
	}
}

func Build() {
	path := tool.GetCurrentPath()
	appJson := path + "/tiny.json"
	if open, err := os.Open(appJson); err == nil {
		open.Close()
	} else if os.IsNotExist(err) {
		fmt.Println("There is no project in the current path")
		return
	}
	os.Setenv("ANDROID_BUILD_CONFIG", appJson)
	androidDir := os.Getenv("ELEPHANT_DIR")
	if androidDir == "" {
		androidDir = path + "/android/"
	}
	tool.ExecCmd(androidDir+"/gradlew", "assembleDebug", "-p", androidDir)

	tool.Adb("install", "-r", androidDir+"build/outputs/apk/debug/app-debug.apk")

	buildConfig := tool.DeCodeAppJson(appJson)

	applicationId := tool.GetApplicationId(buildConfig)

	tool.Adb("shell", "am", "start", "-n", applicationId+"/com.sunmi.android.elephant.core.splash.SplashActivity")
}
