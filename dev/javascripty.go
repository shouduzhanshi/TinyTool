package dev

import (
	"github.com/pterm/pterm"
	"io/ioutil"
	"tiny_tool/build"
	"tiny_tool/observer"
	"tiny_tool/server"
	"tiny_tool/tool"
)

func BuildByJavaScript() {

	go server.StartServer()

	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("building ...")

	config := build.CreateAndroidBuildConfig(tool.GetCurrentPath() + "/src",nil)

	build.AndroidDebug(installApk, func(strings []string) {
		for _, str := range strings {
			pterm.Error.Println(str)
		}
	}, config)

	introSpinner.Stop()

	observer.MonitorSrc(func(js string) {
		m := make(map[string]interface{})
		m["type"] = "changeFile"
		if data, err := ioutil.ReadFile(js); err == nil {
			m["data"] = string(data)
			server.PublishMsg(m)
		}
	})
}
