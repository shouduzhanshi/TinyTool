package build

import (
	"MockConfig/module"
	"MockConfig/tool"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var isInstall = false

var channel = make(chan int)

func ByES6(projectPath, appJsonPath string, appConfig *module.BuildConfig) {
	if _, result, err := tool.BaseCmd("npm", false, "run", "build", "--prefix", projectPath+"/webpack"); err != nil {
		panic(err)
	} else {
		if isBuildJsSuccess(result) {
			startUp(projectPath, *appConfig)
		} else {
			return
		}
	}
	go monitorSrc(projectPath, appJsonPath, onJsFileChange)
	channel <- 1
}

func onJsFileChange(projectPath, path string) {
	fmt.Println("change File:", path)
	if watch, err := fsnotify.NewWatcher(); err != nil {
		panic(err)
	} else {
		watch.Add(projectPath + "/build")
		changeFile := make([]string, 0)
		ints := make(chan int)
		go func(watcher *fsnotify.Watcher) {
			for {
				select {
				case ev := <-watch.Events:
					{
						if ev.Op&fsnotify.Write == fsnotify.Write {
							name := ev.Name
							if strings.HasSuffix(name, ".js") {
								changeFile = append(changeFile, name)
							}
						}
					}
				case <-ints:
					{
						watcher.Close()
						return
					}
				case err := <-watch.Errors:
					fmt.Println(err)
				}
			}
		}(watch)
		tool.ExecCmd("npm", "run", "build", "--prefix", projectPath+"/webpack")
		ints <- 1
		fmt.Println(changeFile)
		close(ints)

	}
}

func monitorSrc(projectPath, appJsonPath string, callback func(string, string)) {
	if watch, err := fsnotify.NewWatcher(); err != nil {
		panic(err)
	} else {
		defer watch.Close()
		dir := tool.GetAllDir(projectPath + "/src")
		//watch.Add(appJsonPath)
		for _, value := range dir {
			watch.Add(value)
		}
		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Write == fsnotify.Write {
						name := ev.Name
						if strings.HasSuffix(name, ".js") {
							callback(projectPath, name)
						} else if strings.HasSuffix(name, ".json") {
							fmt.Println(name)
						}
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err)
				}
			}
		}
	}
}

func isBuildJsSuccess(result []string) bool {
	for _, value := range result {
		if strings.Contains(value, "compiled successfully") {
			return true
		}
	}
	return false
}

func startUp(projectPath string, appConfig module.BuildConfig) {
	isInstall = true
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
	fmt.Println(androidDir)
	buildAndroid(androidDir)
	online := tool.DeviceOnline()
	if online != nil {
		tool.Adb("-s", online.Id, "install", "-r", androidDir+"/build/outputs/apk/debug/app-debug.apk")
		tool.Adb("-s", online.Id, "shell", "am", "start", "-n", appConfig.Build.ApplicationId+"/com.sunmi.android.elephant.core.splash.SplashActivity")
	} else {
		fmt.Println("device not found!")
	}

}
