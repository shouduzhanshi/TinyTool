package log

import (
	"github.com/pterm/pterm"
)

var primaryGreen = pterm.NewStyle(pterm.FgLightGreen, pterm.Bold)

var primaryRed = pterm.NewStyle(pterm.FgRed, pterm.Bold)

func init() {
	pterm.PrintDebugMessages = true

}

func E(data ...interface{}) {
	if len(data) <= 0 {
		return
	}
	primaryRed.Println(data...)
	//primaryGreen.Print("> ")
}

func V(data ...interface{}) {
	if len(data) <= 0 {
		return
	}
	primaryGreen.Println(data...)
	//primaryGreen.Print("> ")
}

func Console(level float64, raw ...interface{}) {
	if len(raw) <= 0 {
		return
	}
	if level == 3.0 {
		pterm.Debug.Println(raw...)
		//primaryGreen.Print("> ")
	}  else if level == 5.0 {
		pterm.Warning.Println(raw...)
		//primaryGreen.Print("> ")
	} else if level == 6.0 {
		pterm.Error.Println(raw...)
		//primaryGreen.Print("> ")
	}else{
		pterm.Info.Println(raw...)
		//primaryGreen.Print("> ")
	}

}

type EchoLogger struct {
}

func (EchoLogger) Write(p []byte) (n int, err error) {
	return len(p), nil
}



func Clean() {
	print("\033[H\033[2J")
}

func Header()  {
	pterm.DefaultHeader.WithFullWidth().Println("Welcome use TinyUI")
}
