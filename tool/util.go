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
	cmd, _, err := BaseCmd("adb", true, raw...)
	return cmd, err
}

func ExecCmd(shell string, raw ...string) (int, error) {
	cmd, _, err := BaseCmd(shell, false, raw...)
	return cmd, err
}

func BaseCmd(shell string, mute bool, raw ...string) (int, []string, error) {
	cmd := exec.Command(shell, raw...)
	defer func() {
		cmd.Process.Release()
		cmd.Process.Kill()
	}()
	result := make([]string, 0)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Log(err.Error())
		return 0, result, nil
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Log(err.Error())
		return 0, result, nil
	}
	if err := cmd.Start(); err != nil {
		Log(err.Error())
		return 0, result, nil
	}
	s := bufio.NewScanner(io.MultiReader(stdout, stderr))

	for s.Scan() {
		text := s.Text()
		result = append(result, text)
		if !mute {
			fmt.Println(text)
		}
	}
	if err := cmd.Wait(); err != nil {
		Log(err.Error())
	}
	return cmd.ProcessState.ExitCode(), result, nil
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

func Cmd(shell string, channel chan int, raw ...string) *os.Process {
	cmd := exec.Command(shell, raw...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Log(err.Error())
		return nil
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Log(err.Error())
		return nil
	}
	if err := cmd.Start(); err != nil {
		Log(err.Error())
		return nil
	}
	go readInput(stdout, stderr, cmd, channel)
	return cmd.Process
}

func readInput(stdout io.ReadCloser, stderr io.ReadCloser, cmd *exec.Cmd, channel chan int) {
	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for s.Scan() {
		text := s.Text()
		if strings.Contains(text, "compiled successfully") {
			channel <- 200
		} else {
			fmt.Println(text)
		}
	}
	fmt.Println("done1")
	if err := cmd.Wait(); err != nil {
		Log(err.Error())
		//channel <- -1
	}
	fmt.Println("done")
}

func GetAllDir(projectPath string) []string {
	dirs := make([]string, 0)
	file, _ := ioutil.ReadDir(projectPath)
	for _, info := range file {
		if info.IsDir() {
			childDir := projectPath + "/" + info.Name()
			dirs = append(dirs, childDir)
			dir := GetAllDir(childDir)
			if len(dir) > 0 {
				dirs = append(dirs, dir...)
			}
		}
	}
	return dirs
}

func DeviceOnline() *module.Device {
	list := GetDeviceList()
	for i := 0; i < len(list); i++ {
		if list[i].Online {
			return &list[i]
		}
	}
	return nil
}

func GetDeviceList() []module.Device {
	deviceList := make([]module.Device, 0)
	if _, result, err := BaseCmd("adb", true, "devices"); err != nil {
		panic(err)
	} else {
		if len(result) > 2 {
			result = result[1 : len(result)-1]
			for i := 0; i < len(result); i++ {
				s := strings.Split(result[i], "\t")
				deviceList = append(deviceList, module.Device{
					Id:     s[0],
					Online: s[1] == "device",
				})
			}
		}
	}
	return deviceList
}
