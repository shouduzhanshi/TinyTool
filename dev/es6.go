package dev

import (
	"bytes"
	"container/list"
	"github.com/pterm/pterm"
	"io/ioutil"
	"net/http"
	"os"
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
		config := tool.GetAppConfig()
		for {
			time.Sleep(time.Duration(500) * time.Millisecond)
			if resp, err := http.Get("http://127.0.0.1:1323/qrCode"); err == nil {
				if resp.StatusCode == 200 {
					data := make(map[string]interface{})
					data["type"] = "apk"
					data["url"] = server.GetApkDownloadUrl()
					server.PublishMsg(data, 0)
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

func ByES6() {
	go server.StartServer()

	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("building ...")

	build.Webpack(buildApk, func(err []string) {
		panic("dsl build fail")
	})

	introSpinner.Stop()

	log.Clean()

	observer.MonitorSrc(onJsFileChange)
}

func onJsFileChange(js string) {
	start := time.Now().UnixNano()
	observer.OnJSFileChange(js, func(list *list.List) {
		sendChangeFile(list, start, js)
	})
}

func sendChangeFile(changeFile *list.List, start int64, js string) {
	config := tool.GetAppConfig()
	pages := make([]module.HotReloadModule, 0)
	editing := false
	for i := changeFile.Front(); i != nil; i = i.Next() {
		s := i.Value.(string)
		if open, err := os.Open(s); err == nil {
			defer open.Close()
			if stat, err := open.Stat(); err == nil {
				name := stat.Name()
				for _, page := range config.Runtime.Pages {
					if page.Name == name[0:len(name)-3] {
						if tool.GetAbsPath(tool.GetCurrentPath(), page.Source) == js {
							editing = true
						}
						if data, err := ioutil.ReadAll(open); err == nil {
							pages = append(pages, module.HotReloadModule{
								Name:    name,
								Router:  page.Router,
								Data:    bytes.NewBuffer(data).String(),
								Editing: editing,
							})
						}
						break
					}
				}
			}
		}
	}
	m := make(map[string]interface{})
	m["type"] = "changeFiles"
	m["files"] = pages
	if !editing {
		m["launcherRouter"] = config.Runtime.LauncherRouter
	}
	server.PublishMsg(m, start)
}
