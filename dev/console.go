package dev

import (
	"github.com/lmorg/readline"
	"math"
	"strings"
	"tiny_tool/server"
)

var items = []any{}

func init() {
	server.UpdateKeys = func(i []interface{}) {
		items = i
		//fmt.Println(i)
		//for _, value := range i {
		//	items = append(i,value.(string))
		//}
		//fmt.Println("update dep")
	}

}

func Console(onExit func()) {
	// Create a new readline instance
	rl := readline.NewInstance()
	rl.MaxTabCompleterRows = math.MaxInt
	//rl.MaxTabItemLength = math.MaxInt
	//rl.MaxTabCompleterRows = math.MaxInt
	// Attach the tab-completion handler (function defined below)
	rl.TabCompleter = Tab
	for {
		// Call readline - which will put the terminal into a pseudo-raw mode
		// and then read from STDIN. After the user has hit <ENTER> the terminal
		// is put back to it's original mode.
		//
		// In this example, `line` is a returned string of the key presses
		// typed into readline.
		line, err := rl.Readline()
		if line != "" && err == nil {
			m := make(map[string]interface{})
			m["type"] = "fragment"
			m["source"] = line
			server.PublishMsg(m,1)
		} else if err != nil && err.Error()=="Ctrl+C" {
			onExit()
			return
		}
	}
}

// items is an example list of possible suggestions to display in readline's
// tab-completion. For the perpose of this example, I basically just grabbed
// a few entries from some random dictionary of terms.

// Tab is the tab-completion handler for this readline example program
func Tab(line []rune, pos int, dtx readline.DelayedTabContext) (string, []string, map[string]string, readline.TabDisplayType) {
	var suggestions []string
	for i := range items {
		s := items[i].(string)
		if strings.HasPrefix(s, string(line)) {
			suggestions = append(suggestions, s[pos:])
		}
	}
	return string(line[:pos]), suggestions, nil, readline.TabDisplayGrid
}
