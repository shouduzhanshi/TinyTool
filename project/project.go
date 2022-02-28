package project

import (
	"MockConfig/tool"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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
	m := make(map[string]interface{})
	json.NewDecoder(bytes.NewBufferString(appConfig)).Decode(&m)
	m2 := m["build"].(map[string]interface{})
	m4 := m["runtime"].(map[string]interface{})
	m4["launcherRouter"] = "hello"
	pages := m4["pages"].([]interface{})
	m5 := make(map[string]interface{})
	m5["router"] = "hello"
	m5["name"] = "hello"
	m5["source"] = "./src/pages/hello/index.js"
	m4["pages"] = append(pages, m5)
	m3 := m2["appName"].(map[string]interface{})
	m3["default"] = projectName
	if marshal, err := json.Marshal(m); err == nil {
		appJson := projectPath + "/tiny.json"
		if create, err := os.Create(appJson); err == nil {
			create.Write(marshal)
			create.Sync()
			create.Close()
		}
	} else {
		panic(err)
	}
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

var appConfig = `
{
"runtime": {
"baseWidth": 750,
"launcherRouter": "",
"pages": [
]
},
"build": {
"appName": {
"default": ""
},
"applicationId": "com.sunmi.elephant.demo",
"versionCode": 1,
"versionName": "0.0.1",
"keystore": {
"storeFilePath": "./android/signing/debug.keystore",
"storePassword": "com.sunmi.elephant.debug",
"keyAlias": "debug",
"keyPassword": "com.sunmi.elephant.debug"
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
"dependencies": [

]
}
}
`
