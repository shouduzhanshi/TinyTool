package dev

import (
	"github.com/pterm/pterm"
	"net/http"
	"time"
	"tiny_tool/build"
	"tiny_tool/log"
	"tiny_tool/observer"
	"tiny_tool/server"
	"tiny_tool/tool"
)

var start = time.Now().Unix()

func buildApk() {
	build.AndroidDebug(installApk, func(err []string) {
		for _, s := range err {
			pterm.Error.Println(s)
		}
	}, build.CreateAndroidBuildConfig(tool.GetCurrentPath()+"/build",nil))
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
	log.E("total duration ", time.Now().Unix()-start, " s")
}

func ByES6() {

	go server.StartServer()

	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("building ...")

	build.Webpack(buildApk, func(err error) {
		panic(err)
	})

	introSpinner.Stop()

	log.Clean()

	observer.MonitorSrc(observer.OnJSFileChange)
}
