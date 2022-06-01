package dev

import (
	"bytes"
	"encoding/json"
	"github.com/martinlindhe/notify"
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
		log.E("All devices are offline, please check the device usb connection")
		notify.Notify("Tiny CLI", "warning", "All devices are offline, please check the device usb connection", "")
		config := tool.GetAppConfig()
		for {
			time.Sleep(time.Duration(500) * time.Millisecond)
			if resp, err := http.Get("http://127.0.0.1:1323/qrCode"); err == nil {
				if resp.StatusCode == 200 {
					data := make(map[string]interface{})
					data["type"] = "apk"
					data["url"] = server.GetApkDownloadUrl()
					server.PublishMsg(data, 1)
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
	cacheChangeFile := changeFile
	changeFile = make([]string, 0)
	pages := make([]module.HotReloadModule, 0)
	var size int64
	for _, js := range cacheChangeFile {
		config := tool.GetAppConfig()
		if open, err := os.Open(js); err == nil {
			defer open.Close()
			if stat, err := open.Stat(); err == nil {
				name := stat.Name()
				size += stat.Size()
				for _, page := range config.Runtime.Pages {
					if page.Name == name[0:len(name)-3] {
						pages = append(pages, module.HotReloadModule{
							Name:     page.Name,
							Router:   page.Router,
							Data:     server.GetServerUrl() + "build/" + name,
							Size:     stat.Size(),
							FileName: name,
						})
						break
					}
				}
			}
		}
	}
	if len(pages) > 0 {
		m := make(map[string]interface{})
		m["type"] = "changeFiles"
		m["files"] = pages
		m["size"] = size
		server.PublishMsg(m, 1)
		log.E("size ", size, "bytes")
	} else {
		log.E("unchanged !")
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
