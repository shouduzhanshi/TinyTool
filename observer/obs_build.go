package observer

import (
	"container/list"
	"github.com/fsnotify/fsnotify"
	"github.com/pterm/pterm"
	"io/ioutil"
	"strings"
	"time"
	"tiny_tool/build"
	"tiny_tool/log"
	"tiny_tool/server"
	"tiny_tool/tool"
)

func OnJSFileChange(string) {
	start := time.Now().UnixNano()
	var closeWatchChannel = make(chan int)
	if watch, err := fsnotify.NewWatcher(); err == nil {
		projectPath := tool.GetCurrentPath()
		watch.Add(projectPath + "/dev")
		changeFile := list.New()
		go buildDirChangeCallback(watch, closeWatchChannel, changeFile)
		introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("building ...")
		build.Webpack(func() {
			go sendChangeFile(changeFile, start)
		}, func(err error) {

		})
		introSpinner.Stop()
		closeWatchChannel <- 1
		close(closeWatchChannel)
	} else {
		panic(err)
	}
}

func buildDirChangeCallback(watcher *fsnotify.Watcher, closeWatchChannel chan int, changeFile *list.List) {
	for {
		select {
		case ev := <-watcher.Events:
			{
				if ev.Op&fsnotify.Write == fsnotify.Write {
					name := ev.Name
					if strings.HasSuffix(name, ".js") {
						changeFile.PushBack(name)
					}
				}
			}
		case <-closeWatchChannel:
			{
				watcher.Close()
				return
			}
		case err := <-watcher.Errors:
			log.E("watcher.Errors", err)
		}
	}
}

func sendChangeFile(changeFile *list.List, start int64) {
	for i := changeFile.Front(); i != nil; i = i.Next() {
		go sending(i.Value.(string), start)
	}
}

func sending(changeFile string, start int64) {
	m := make(map[string]interface{})
	m["type"] = "changeFile"
	if data, err := ioutil.ReadFile(changeFile); err == nil {
		m["data"] = string(data)
		server.PublishMsg(m, start)
	}
}
