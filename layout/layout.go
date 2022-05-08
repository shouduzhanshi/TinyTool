package layout

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"io/ioutil"
	"os"
	"time"
	"tiny_tool/log"
	"tiny_tool/tool"
)

func PrintLayout(data string) {
	//if true {
	//	return
	//}
	layoutData := AndroidLayoutData{}
	//{Level: 0, Text: "Level 0"},
	//{Level: 1, Text: "Level 1"},
	//{Level: 2, Text: "Level 2"},
	list := list.New()
	json.Unmarshal(bytes.NewBufferString(data).Bytes(), &layoutData)
	layout := layoutData.Data
	level := 0
	funcName(level, layout, list)
	items := make([]pterm.LeveledListItem, 0)
	for i := list.Front(); i != nil; i = i.Next() {
		items = append(items, i.Value.(pterm.LeveledListItem))
	}
	log.Clean()
	pterm.DefaultTree.WithRoot(pterm.NewTreeFromLeveledList(items)).Render()
	str,_ := json.Marshal(layout)
	s := tool.GetCurrentPath() + "/" + fmt.Sprint(time.Now().UnixNano()) + ".json"
	ioutil.WriteFile(s,str,os.ModePerm.Perm())
}

func funcName(level int, layout AndroidLayout, items *list.List) {

	for _, data := range layout.AndroidNodeList {
		items.PushBack(pterm.LeveledListItem{
			Level: level,
			//+" x"+strconv.(layout.X) +" y"+layout.Y +" width"+layout.Width+" height"+layout.Height
			Text: data.LayoutName + " "+ fmt.Sprint(data.Info),
		})
		if len(data.AndroidNodeList) > 0 {

			funcName(level+1, data, items)
		}
	}
}

type AndroidLayoutData struct {
	Type string
	Data AndroidLayout
}

type AndroidLayout struct {
	Width           int
	Height          int
	X               float32
	Y               float32
	LayoutName      string
	AndroidNodeList []AndroidLayout
	Info            map[string]interface{}
}
