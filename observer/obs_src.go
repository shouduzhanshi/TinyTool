package observer

import (
	"github.com/fsnotify/fsnotify"
	"strings"
	"tiny_tool/log"
	"tiny_tool/tool"
)

func MonitorSrc(srcPath string, callback func()) {
	if watch, err := fsnotify.NewWatcher(); err != nil {
		panic(err)
	} else {
		defer watch.Close()
		dir := tool.GetAllDir(srcPath)
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
							callback()
						}
					}
				}
			case err := <-watch.Errors:
				{
					log.LogV("error : ", err)
				}
			}
		}
	}
}
