package install

import (
	"MockConfig/tool"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func Install() {
	homePath := os.Getenv("HOME")
	s := homePath + "/tiny"
	saveFile := s + "/tiny_tool"
	if file, err := os.Open(saveFile); err == nil {
		file.Close()
		fmt.Println("success")
		return
	} else if os.IsNotExist(err) {
		binPath := tool.GetCurrentPath() + "/tiny_tool"
		if open, err := os.Open(s); err == nil {
			open.Close()
		} else if os.IsNotExist(err) {
			if err := os.Mkdir(s, os.ModePerm); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
		if create, err := os.Create(saveFile); err == nil {
			defer create.Close()
			if open, err := os.Open(binPath); err == nil {
				defer open.Close()
				if written, err := io.Copy(create, open); err == nil {
					fmt.Println("copy:", written)
					os.Chmod(saveFile, 755)
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
		os.Remove(binPath)
		if os.Getenv("SHELL") == "/bin/zsh" {
			setEnv(homePath, "/.zshrc")
		} else if os.Getenv("SHELl") == "bin/bsh" {
			setEnv(homePath, "/.bshrc")
		}
		fmt.Println("install success")
	}

}

func setEnv(homePath, envSource string) {
	if create, err := os.OpenFile(homePath+envSource,os.O_RDWR, 0); err == nil {
		defer create.Close()
		scaner := bufio.NewScanner(create)
		isSetEnv := false
		for scaner.Scan() {
			line := scaner.Text()
			if strings.Contains(line, "TINY_TOOL") {
				isSetEnv = true
				break
			}
		}
		if !isSetEnv {
			create.WriteString("\n")
			writeEnv(homePath, create)
		}
	} else if create, err := os.Create(homePath + envSource); err == nil {
		writeEnv(homePath, create)
	}else {
		fmt.Println(err)
	}
}

func writeEnv(homePath string, create *os.File) {
	create.WriteString("export TINY_TOOL=~/tiny")
	create.WriteString("\n")
	create.WriteString("export PATH=$TINY_TOOL:$PATH")
	create.Sync()
}
