package build

import (
	"MockConfig/log"
	"MockConfig/module"
	"MockConfig/observer"
	"MockConfig/server"
	"MockConfig/tool"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func ByES6(projectPath string, appConfig *module.BuildConfig) {
	start := time.Now().Unix()
	if _, result, err := tool.BaseCmd("npm", false, "run", "build", "--prefix", projectPath+"/webpack"); err == nil {
		log.LogE("npm build duration ", time.Now().Unix()-start, " s")
		if isBuildJsSuccess(result) {
			startUp(projectPath, *appConfig, start)
		} else {
			return
		}
	} else {
		panic(err)
	}
	go observer.MonitorSrc(projectPath+"/src", observer.OnJSFileChange)
	server.StartServer()
}

func isBuildJsSuccess(result []string) bool {
	for _, value := range result {
		if strings.Contains(value, "compiled successfully") {
			return true
		}
	}
	return false
}

func startUp(projectPath string, appConfig module.BuildConfig, start int64) {
	dslDir := projectPath + "/build"
	pages := make([]module.PageModule, 0)
	if files, err := ioutil.ReadDir(dslDir); err == nil {
		for _, file := range files {
			fileName := file.Name()
			if strings.HasSuffix(fileName, ".js") {
				modules := appConfig.Runtime.Pages
				for _, page := range modules {
					if page.Name == fileName[0:len(fileName)-3] {
						page.Source = dslDir + "/" + fileName
						pages = append(pages, page)
					}
				}
			}
		}
	}
	appConfig.Runtime.Pages = pages
	appConfig.Build.Keystore.StoreFilePath = tool.GetAbsPath(projectPath, appConfig.Build.Keystore.StoreFilePath)
	icon := appConfig.Build.LauncherIcon

	for i := 0; i < len(icon); i++ {
		icon[i].Icon = tool.GetAbsPath(projectPath, icon[i].Icon)
	}

	splash := appConfig.Build.Splash.Background

	for i := 0; i < len(splash); i++ {
		splash[i].Src = tool.GetAbsPath(projectPath, splash[i].Src)
	}

	appJson := projectPath + "/build/.mock.json"

	appConfig.Runtime.Ws = server.GetWsPath()

	if data, err := json.Marshal(appConfig); err != nil {
		panic(err)
	} else {
		if err := ioutil.WriteFile(appJson, data, os.ModePerm); err != nil {
			panic(err)
		}
	}

	defer os.Remove(appJson)

	os.Setenv("ANDROID_BUILD_CONFIG", appJson)

	androidDir := projectPath + "/android"

	if len(os.Args) > 2 {
		androidDir = os.Args[2]
	}
	androidBuildDuration := time.Now().Unix()
	buildAndroid(androidDir)
	log.LogE("android build duration ", time.Now().Unix()-androidBuildDuration, " s")
	list := tool.GetDeviceList()
	for _, device := range list {
		if device.Online {
			log.LogV("install app to ", device.Id, " ....")
			installStart := time.Now().Unix()
			tool.Adb("-s", device.Id, "install", "-r", androidDir+"/build/outputs/apk/debug/app-debug.apk")
			log.LogV("install app to ", device.Id, " duration ", time.Now().Unix()-installStart, " s")
			openStart := time.Now().Unix()
			tool.Adb("-s", device.Id, "shell", "am", "start", "-n", appConfig.Build.ApplicationId+"/com.sunmi.android.elephant.core.splash.SplashActivity")
			log.LogV("open app from ", device.Id," ", time.Now().Unix()-openStart, " s ")
		}
	}
	log.LogE("total duration ", time.Now().Unix()-start, " s")
	if tool.DeviceOnline() == nil {
		log.LogE("device not found!")
		if !appConfig.DisableOpenBrowser {
			go openBrowser()
		}
	}

}

func progressBar(isInstallSuccess *bool, title, doneMsg string) {
	go func(isInstallSuccess *bool) {
		fmt.Print(title + " .....")
		for {
			time.Sleep(time.Duration(200) * time.Millisecond)
			if !*isInstallSuccess {
				fmt.Print(".")
			} else {
				fmt.Print(doneMsg, "\r\n")
				return
			}
		}
	}(isInstallSuccess)
}

func openBrowser() {
	for {
		time.Sleep(time.Duration(500) * time.Millisecond)
		if resp, err := http.Get("http://127.0.0.1:1323/qrCode"); err == nil {
			if resp.StatusCode == 200 {
				tool.ExecCmd("open", "-a", "Google Chrome", "http://127.0.0.1:1323")
				return
			}
			resp.Body.Close()
		}
	}
}
