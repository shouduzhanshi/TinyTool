package project

import (
	"MockConfig/module"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func simpleProject(projectName string, projectPath string) {

	writeGitignore(projectPath,"*/.DS_Store\r\n.idea\r\nnode_modules\r\ndist\r\nandroid")

	writeAppConfig(projectName, projectPath,module.JavaScript)

	mkSRCdir(projectPath)

	fmt.Println("success")
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
}

func writeAppConfig(projectName , projectPath ,projectType string) {
	appJson := projectPath + "/tiny.json"
	appConfig = strings.ReplaceAll(appConfig, "$PROJECT_NAME", projectName)
	appConfig = strings.ReplaceAll(appConfig, "$APPID", "com.sunmi.elephant."+strings.ToLower(projectName))
	appConfig = strings.ReplaceAll(appConfig, "$PROJECT_TYPE", projectType)
	ioutil.WriteFile(appJson, bytes.NewBufferString(appConfig).Bytes(), os.ModePerm)
}

func writeGitignore(projectPath, gitignore string) {
	if file, err := os.Create(projectPath + "/.gitignore"); err == nil {
		file.WriteString(gitignore)
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
  },
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
  "projectType": "$PROJECT_TYPE"
}
`
