package dev

import (
	"bytes"
	"encoding/json"
	"github.com/pterm/pterm"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"
	"time"
	"tiny_tool/build"
	"tiny_tool/log"
	"tiny_tool/module"
	"tiny_tool/observer"
	"tiny_tool/server"
	"tiny_tool/tool"
)

var start = time.Now().UnixNano()

func buildApk() {
	build.AndroidDebug(installApk, func(err []string) {
		for _, s := range err {
			pterm.Error.Println(s)
		}
	}, build.CreateAndroidBuildConfig(tool.GetCurrentPath()+"/build", nil))
}

func installApk() {
	InstallApk(nil, func() {
		log.E("No active device found")
		config := tool.GetAppConfig()
		for {
			time.Sleep(time.Duration(500) * time.Millisecond)
			if resp, err := http.Get("http://127.0.0.1:1323/qrCode"); err == nil {
				if resp.StatusCode == 200 {
					data := make(map[string]interface{})
					data["type"] = "apk"
					data["url"] = server.GetApkDownloadUrl()
					server.PublishMsg(data)
					if !config.DisableOpenBrowser {
						tool.ExecCmd("open", "-a", "Google Chrome", "http://127.0.0.1:1323")
					}
					resp.Body.Close()
					return
				}
				resp.Body.Close()
			}
		}
	})
	log.E("total duration ", (time.Now().UnixNano()-start)/1e6, " ms")
}

var isInitDone = false

func ByES6() {
	go server.StartServer()
	go Console(func() {
		if cmd != nil {
			cmd.Process.Signal(syscall.SIGINT)
			cmd = nil
		}
		os.Exit(0)
	})

	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("building ...")

	WebpackDev(func() {
		if !isInitDone {
			isInitDone = true
			buildApk()
		} else {
			buildSuccess()
		}
	})

	introSpinner.Stop()

	log.Clean()

	observer.Dev(onJsChange, onConfigChange)
}

func buildSuccess() {

	pages := make([]module.HotReloadModule, 0)
	for _,js := range changeFile {
		config := tool.GetAppConfig()
		if open, err := os.Open(js); err == nil {
			defer open.Close()
			if stat, err := open.Stat(); err == nil {
				name := stat.Name()
				for _, page := range config.Runtime.Pages {
					if page.Name == name[0:len(name)-3] {
						if data, err := ioutil.ReadAll(open); err == nil {
							pages = append(pages, module.HotReloadModule{
								Name:   name,
								Router: page.Router,
								Data:   bytes.NewBuffer(data).String(),
							})
						}
						break
					}
				}
			}
		}
	}
	if len(pages)>0 {
		m := make(map[string]interface{})
		m["type"] = "changeFiles"
		m["files"] = pages
		server.PublishMsg(m)
		for _, page := range pages {
			log.E("change page"+page.Router)
		}
		changeFile = make([]string, 0)
	}else{
		log.E("no change!")
	}
}

func onConfigChange(path string) {
	if file, err := ioutil.ReadFile(path); err == nil {
		decoder := json.NewDecoder(bytes.NewBuffer(file))
		buildConfig := module.BuildConfig{}
		if err := decoder.Decode(&buildConfig); err == nil {
			cmd.Process.Signal(syscall.SIGINT)
			WebpackDev(func() {
				buildSuccess()
			})
		}
	}
}

var changeFile = make([]string, 0)

func onJsChange(js string) {
	if !isInitDone {
		return
	}
	changeFile = append(changeFile, js)
	log.E(js)
}
