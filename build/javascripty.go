package build

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"tiny_tool/module"
	"tiny_tool/tool"
)


func HotReloadByJavaScript() {
	path := tool.GetCurrentPath()
	appJson := path + "/"+module.TINY_JSON
	json := tool.DeCodeAppJson(appJson)
	applicationId := tool.GetApplicationId(*json)
	if json.ProjectType != module.JavaScript {
		return
	}
	if adb, _ := tool.Adb("shell", "pm", "path", applicationId); adb == 1 {
		BuildByJavaScript(true)
		return
	}
	changeFilePath := os.Args[2]
	tool.Adb("shell", "mkdir", "sdcard/Android/data/"+applicationId)
	tool.Adb("shell", "mkdir", "sdcard/Android/data/"+applicationId+"/cache")
	tool.Adb("push", changeFilePath, "sdcard/Android/data/"+applicationId+"/cache")
	if file, err := os.Open(changeFilePath); err == nil {
		defer file.Close()
		if info, err := file.Stat(); err == nil {
			tool.Adb("shell", "am", "start", "-n", applicationId+"/com.sunmi.android.elephant.api.container.ContainerActivity", "-f", "0x10000000", "--es", "hotReLoad", "\"hotReLoad://cache/"+info.Name()+"\"")
		}
	}
}


func BuildByJavaScript(isHotReload bool) {
	path := tool.GetCurrentPath()
	appJson := path + "/"+module.TINY_JSON
	if open, err := os.Open(appJson); err == nil {
		open.Close()
	} else if os.IsNotExist(err) {
		fmt.Println("There is no project in the current path")
		return
	}

	buildConfig := tool.DeCodeAppJson(appJson)

	androidDir := getElephantDir(isHotReload)
	if androidDir == "" {
		androidDir = path + "/android"
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

	buildAndroid(androidDir)

	installAndroidApp(androidDir)

	applicationId := tool.GetApplicationId(*buildConfig)

	startApp(applicationId)
}
