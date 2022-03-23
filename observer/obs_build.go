package observer

import (
	"container/list"
	"github.com/fsnotify/fsnotify"
	"github.com/pterm/pterm"
	"strings"
	"tiny_tool/build"
	"tiny_tool/log"
	"tiny_tool/tool"
)

func OnJSFileChange(js string,callback func(*list.List)) {
	var closeWatchChannel = make(chan int)
	if watch, err := fsnotify.NewWatcher(); err == nil {
		projectPath := tool.GetCurrentPath()
		watch.Add(projectPath + "/build")
		changeFile := list.New()
		go buildDirChangeCallback(watch, closeWatchChannel, changeFile)
		introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("building ...")
		build.Webpack(func() {
			go callback(changeFile)
		}, func(err []string) {
			panic("dsl build fail")
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
