package build

import (
	"encoding/json"
	"github.com/pterm/pterm"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"tiny_tool/log"
	"tiny_tool/module"
	"tiny_tool/observer"
	"tiny_tool/server"
	"tiny_tool/tool"
)

func ByES6(projectPath string, appConfig *module.BuildConfig) {
	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("building ...")
	start := time.Now().Unix()
	if _, result, err := tool.BaseCmd("npm", false,"run", "build", "--prefix", projectPath+"/webpack"); err == nil {
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
	introSpinner.Stop()
	log.Clean()
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
	android := buildAndroid(androidDir)
	if android!=0 {
		panic("android build fail~")
		return
	}
	log.LogE("build android state ",android)
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
