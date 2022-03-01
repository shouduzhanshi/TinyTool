package tool

import (
	"MockConfig/module"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func Log(msg string) {
	fmt.Println(msg)
}

func Adb(raw ...string) (int, error) {
	return ExecCmd("adb", raw...)
}

func ExecCmd(shell string, raw ...string) (int, error) {
	cmd := exec.Command(shell, raw...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Log(err.Error())
		return 0, nil
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Log(err.Error())
		return 0, nil
	}
	if err := cmd.Start(); err != nil {
		Log(err.Error())
		return 0, nil
	}
	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for s.Scan() {
		text := s.Text()
		Log(text)
	}
	if err := cmd.Wait(); err != nil {
		Log(err.Error())
	}

	return cmd.ProcessState.ExitCode(), nil
}

func DeCodeAppJson(appJson string) *module.BuildConfig {
	if file, err := ioutil.ReadFile(appJson); err == nil {
		decoder := json.NewDecoder(bytes.NewBuffer(file))
		buildConfig := module.BuildConfig{}
		decoder.Decode(&buildConfig)
		return &buildConfig
	} else {
		panic(err)
	}
	return nil
}

func GetApplicationId(buildConfig module.BuildConfig) string {
	return buildConfig.Build.ApplicationId
}

func GetAbsPath(parentPath, absPath string) string {
	if strings.HasPrefix(absPath, "./") {
		return parentPath + "/" + absPath[2:]
	} else if strings.HasPrefix(absPath, "../") {
		runes := []rune(absPath)
		count := 0
		lastIndex := 0
		for i := 0; i < len(runes); i++ {
			if '.' == runes[i] {
				if i+1 < len(runes) && '.' == runes[i+1] {
					if i+2 < len(runes) && '/' == runes[i+2] {
						count++
						lastIndex = i + 3
						i = i + 2
					}
				}
			}
		}
		for count > 0 {
			count--
			parentPath = getParentDirectory(parentPath)
		}
		return parentPath + "/" + absPath[lastIndex:]
	} else {
		return absPath
	}
}




func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}


func GetCurrentPath() string {
	if getwd, err := os.Getwd(); err == nil {
		return getwd
	} else {
		panic(err)
	}
}
