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
	go Console()

	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("building ...")

	WebpackDev(func() {
		isInitDone = true
		buildApk()
	})

	introSpinner.Stop()

	log.Clean()

	observer.Dev(onJsChange, onConfigChange)
}

func onConfigChange(path string) {
	if file, err := ioutil.ReadFile(path); err == nil {
		decoder := json.NewDecoder(bytes.NewBuffer(file))
		buildConfig := module.BuildConfig{}
		if err := decoder.Decode(&buildConfig); err == nil {
			cmd.Process.Signal(syscall.SIGINT)
			WebpackDev(nil)
		}
	}
}

func onJsChange(js string) {
	if !isInitDone {
		return
	}
	config := tool.GetAppConfig()
	pages := make([]module.HotReloadModule, 0)
	if open, err := os.Open(js); err == nil {
		defer open.Close()
		if stat, err := open.Stat(); err == nil {
			name := stat.Name()
			for _, page := range config.Runtime.Pages {
				if page.Name == name[0:len(name)-3] {
					if data, err := ioutil.ReadAll(open); err == nil {
						pages = append(pages, module.HotReloadModule{
							Name:    name,
							Router:  page.Router,
							Data:    bytes.NewBuffer(data).String(),
							Editing: true,
						})
					}
					break
				}
			}
		}
	}
	if len(pages) <= 0 {
		return
	}
	m := make(map[string]interface{})
	m["type"] = "changeFiles"
	m["files"] = pages
	server.PublishMsg(m)
}
