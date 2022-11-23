package build

import (
	"github.com/pterm/pterm"
	"os"
	"strconv"
	"time"
	"tiny_tool/log"
	"tiny_tool/module"
	"tiny_tool/tool"
)

var projectPath = tool.GetCurrentPath()

var androidDir = GetAndroidDir()

func Build() {
	config := tool.GetAppConfig()
	var configPath string
	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("building ...")
	outApk := projectPath + "/" + config.Build.AppName.Default + "_v" + config.Build.VersionName + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".apk"
	if config.ProjectType == module.JavaScript || config.ProjectType == "" {
		configPath = CreateAndroidBuildConfig(tool.GetCurrentPath()+"/src", &outApk)
	} else if config.ProjectType == module.ES6 || config.ProjectType == module.JSX {
		Webpack(func() {

		}, func(err []string) {
			panic("dsl build fail")
		})
		configPath = CreateAndroidBuildConfig(tool.GetCurrentPath()+"/build", &outApk)
	}
	AndroidRelease(func() {
		introSpinner.Stop()
		log.V("build success")
	}, func(strings []string) {
		introSpinner.Stop()
		for _, str := range strings {
			pterm.Error.Println(str)
		}
	}, configPath)
}

func GetAndroidDir() string {
	dir := projectPath + "/android"
	if len(os.Args) > 2 {
		dir = os.Args[2]
	}
	return dir
}



func AndroidDebug(success func(), fail func([]string), appJsonPath string) {
	Android(success, fail, appJsonPath, "assembleDebug")
}

func AndroidRelease(success func(), fail func([]string), appJsonPath string) {
	Android(success, fail, appJsonPath, "assembleRelease")
}

func Android(success func(), fail func([]string), appJsonPath, tag string) {
	defer os.Remove(appJsonPath)
	os.Setenv("ANDROID_BUILD_CONFIG", appJsonPath)
	androidBuildDuration := time.Now().UnixNano()
	mute := false
	if os.Getenv("SHOW_BUILD_LOG") == "false"{
		mute = true
	}
	if cmd, result := tool.BaseCmd(androidDir+"/gradlew", mute, tag, "-p", androidDir); cmd == 0 {
		log.E("android build duration ", (time.Now().UnixNano()-androidBuildDuration)/1e6, " ms")
		success()
	} else {
		fail(result)
		panic("android build fail")
	}
}

func Webpack(success func(), fail func([]string)) {
	start := time.Now().UnixNano()
	mute := false
	if os.Getenv("SHOW_BUILD_LOG") == "false"{
		mute = true
	}
	if code, result := tool.BaseCmd("npm", mute, "run", "build", "--prefix", projectPath); code == 0 {
		log.E("npm build duration ", (time.Now().UnixNano() - start) / 1e6, " ms")
		success()
	} else {
		for _, s := range result {
			pterm.Error.Println(s)
		}
		fail(result)
	}

}
