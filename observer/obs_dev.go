package observer

import (
	"github.com/fsnotify/fsnotify"
	"strings"
	"tiny_tool/log"
	"tiny_tool/tool"
)

func Dev(jsChange func(string), configChange func(string)) {
	if watcher, err := fsnotify.NewWatcher(); err == nil {
		projectPath := tool.GetCurrentPath()
		watcher.Add(projectPath + "/build/")
		watcher.Add(projectPath + "/app.config.json")
		for {
			select {
			case ev := <-watcher.Events:
				{
					if ev.Op&fsnotify.Write == fsnotify.Write {
						name := ev.Name
						if strings.HasSuffix(name, ".js") {
							jsChange(name)
						} else if strings.HasSuffix(name, ".json") {
							configChange(name)
						}
					}
				}
			case err := <-watcher.Errors:
				log.E("watcher.Errors", err)
				panic(err)
			}
		}
	} else {
		panic(err)
	}
}
