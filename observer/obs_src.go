package observer

import (
	"container/list"
	"github.com/fsnotify/fsnotify"
	"os"
	"strings"
	"tiny_tool/log"
	"tiny_tool/tool"
)

var dirs = list.New()

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
					name := ev.Name
					if ev.Op&fsnotify.Write == fsnotify.Write {
						if strings.HasSuffix(name, ".js") {
							callback()
						}
					} else if ev.Op&fsnotify.Create == fsnotify.Create {
						refreshDir(watch, name)
					} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
						clearDir(watch, name)
					}
				}
			case err := <-watch.Errors:
				{
					log.V("error : ", err)
				}
			}
		}
	}
}

func clearDir(watch *fsnotify.Watcher, name string) {
	var dirIndex *list.Element
	for i := dirs.Front(); i != nil; i = i.Next() {
		if i.Value == name {
			dirIndex = i
			break
		}
	}
	if dirIndex != nil {
		dirs.Remove(dirIndex)
		watch.Remove(name)
	}
}

func refreshDir(watch *fsnotify.Watcher, name string) {
	if open, err := os.Open(name); err == nil {
		defer open.Close()
		if stat, err := open.Stat(); err == nil {
			if stat.IsDir() {
				isExits := false
				for i := dirs.Front(); i != nil; i = i.Next() {
					if i.Value == name {
						isExits = true
						break
					}
				}
				if !isExits {
					dirs.PushBack(name)
					watch.Add(name)
				}
			}
		}
	}
}
