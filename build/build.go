package build

import (
	"MockConfig/tool"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func HotReload() {
	path := tool.GetCurrentPath()
	appJson := path + "/tiny.json"
	json := tool.DeCodeAppJson(appJson)
	applicationId := tool.GetApplicationId(*json)
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

	buildConfig := tool.DeCodeAppJson(appJson)

	androidDir := getElephantDir()
	fmt.Println(androidDir)
	if androidDir == "" {
		androidDir = path + "/android/"
	} else {
		buildConfig.Build.Keystore.StoreFilePath = tool.GetAbsPath(path, buildConfig.Build.Keystore.StoreFilePath)

		icon := buildConfig.Build.LauncherIcon
		for i := 0; i < len(icon); i++ {
			icon[i].Icon = tool.GetAbsPath(path, icon[i].Icon)
		}

		splash := buildConfig.Build.Splash.Background
		for i := 0; i < len(splash); i++ {
			splash[i].Src = tool.GetAbsPath(path, splash[i].Src)
		}

		pages := buildConfig.Runtime.Pages
		for i := 0; i < len(pages); i++ {
			pages[i].Source = tool.GetAbsPath(path, pages[i].Source)
		}

		appJson = path + "/.mock.json"
		if marshal, err := json.Marshal(buildConfig); err == nil {
			if err := ioutil.WriteFile(appJson, marshal, os.ModePerm); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}

		defer os.Remove(appJson)
	}

	os.Setenv("ANDROID_BUILD_CONFIG", appJson)

	tool.ExecCmd(androidDir+"/gradlew", "assembleDebug", "-p", androidDir)

	tool.Adb("install", "-r", androidDir+"build/outputs/apk/debug/app-debug.apk")

	applicationId := tool.GetApplicationId(*buildConfig)

	tool.Adb("shell", "am", "start", "-n", applicationId+"/com.sunmi.android.elephant.core.splash.SplashActivity")
}

func getElephantDir() string {
	if len(os.Args) >= 3 {
		return os.Args[2]
	}
	return ""
}
