package observer

import (
	"MockConfig/log"
	"MockConfig/server"
	"MockConfig/tool"
	"container/list"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"strings"
	"time"
)

func OnJSFileChange() {
	start := time.Now().UnixNano()
	var closeWatchChannel = make(chan int)
	if watch, err := fsnotify.NewWatcher(); err == nil {
		projectPath := tool.GetCurrentPath()
		watch.Add(projectPath + "/build")
		changeFile := list.New()
		go buildDirChangeCallback(watch,closeWatchChannel,changeFile)
		npmStart := time.Now().Unix()
		tool.ExecCmd("npm", "run", "build", "--prefix", projectPath+"/webpack")
		end := time.Now().Unix()
		log.LogE("DSL构建耗时: ", end-npmStart, " seconds")
		go sendChangeFile(changeFile, start)
		closeWatchChannel <- 1
		close(closeWatchChannel)
	} else {
		panic(err)
	}
}

func buildDirChangeCallback(watcher *fsnotify.Watcher, closeWatchChannel chan int,changeFile *list.List) {
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
			log.LogE(err)
		}
	}
}


func sendChangeFile(changeFile *list.List,start int64) {
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