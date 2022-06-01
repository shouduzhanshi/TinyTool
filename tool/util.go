package tool

import (
	"io/ioutil"
	"os"
	"strings"
	"tiny_tool/module"
)

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

func GetBuildPath() string {
	return GetCurrentPath() + "/build"
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
	if _, result := BaseCmd("adb", true, "devices"); len(result) > 2 {
		result = result[1 : len(result)-1]
		for i := 0; i < len(result); i++ {
			s := strings.Split(result[i], "\t")
			if len(s) < 2 {
				return deviceList
			}
			deviceList = append(deviceList, module.Device{
				Id:     s[0],
				Online: s[1] == "device",
			})
		}
	}
	return deviceList
}
