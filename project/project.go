package project

import (
	"MockConfig/tool"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func InitProject() {

	projectName := ""

	fmt.Println("Please enter a project name")

	fmt.Scanln(&projectName)

	projectPath := tool.GetCurrentPath()+"/"+projectName

	if err := os.Mkdir(projectPath,os.ModePerm);err !=nil{
		panic(err)
	}

	projectAndroid := projectPath + "/android"

	tool.ExecCmd("git", "clone", "-b", "master", "git@codeup.teambition.com:sunmi/MaxProgram/Android/Elephant.git", projectAndroid)

	writeGitignore(projectPath)

	writeAppConfig(projectName, projectPath)

	mkSRCdir(projectPath)
}

func mkSRCdir(projectPath string) {
	if err := os.Mkdir(projectPath+"/src", os.ModePerm); err == nil {

	} else {
		panic(err)
	}

	if err := os.Mkdir(projectPath+"/src/pages", os.ModePerm); err == nil {

	} else {
		panic(err)
	}

	if err := os.Mkdir(projectPath+"/src/pages/hello", os.ModePerm); err == nil {

	} else {
		panic(err)
	}
	if file, err := os.Create(projectPath + "/src/pages/hello/index.js"); err == nil {
		file.WriteString(demo)
		file.Sync()
		file.Close()
	} else {
		panic(err)
	}
	fmt.Println("success")
}

func writeAppConfig(projectName string, projectPath string) {
	appJson := projectPath + "/tiny.json"
	appConfig = strings.ReplaceAll(appConfig,"$PROJECT_NAME",projectName)
	appConfig = strings.ReplaceAll(appConfig,"$APPID","com.sunmi.elephant."+strings.ToLower(projectName))
	ioutil.WriteFile(appJson,bytes.NewBufferString(appConfig).Bytes(),os.ModePerm)
}

func writeGitignore(projectPath string) {
	if file, err := os.Create(projectPath + "/.gitignore"); err == nil {
		file.WriteString("*/.DS_Store\r\n.idea\r\nnode_modules\r\ndist\r\nandroid")
		file.Sync()
		file.Close()
	} else {
		panic(err)
	}
}

var demo = `
TinyUI.render(
    TinyDOM.createElement("column", {
            style: {
                width:"100%",
                height:"100%"
            },
        }, TinyDOM.createElement("text", {
            style:{
                width:"100%",
                height:"100%",
                textAlign:"center"
            }
        }, "hello word!")
    )
)
`

var appConfig = `{
  "build": {
    "appName": {
      "default": "$PROJECT_NAME"
    },
    "applicationId": "$APPID",
    "dependencies": [],
    "keystore": {
      "keyAlias": "debug",
      "keyPassword": "com.sunmi.elephant.debug",
      "storeFilePath": "./android/signing/debug.keystore",
      "storePassword": "com.sunmi.elephant.debug"
    },
    "launcherIcon": [
      {
        "icon": "./android/mock/ic_launcher.png",
        "resolution": "xxxhdpi"
      }
    ],
    "splash": {
      "background": [
        {
          "resolution": "xxxhdpi",
          "src": "./android/mock/splash.png"
        }
      ]
    },
    "versionCode": 1,
    "versionName": "0.0.1"
  },
  "runtime": {
    "baseWidth": 750,
    "launcherRouter": "hello",
    "pages": [
      {
        "name": "hello",
        "router": "hello",
        "source": "./src/pages/hello/index.js"
      }
    ]
  }
}

`
