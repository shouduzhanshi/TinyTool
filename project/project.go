package project

import (
	"MockConfig/tool"
	"fmt"
	"os"
)

func InitProject() {

	projectName := ""

	projectType := ""
	fmt.Println("Please enter a project name")

	fmt.Scanln(&projectName)

	fmt.Println("Please select project type")

	fmt.Println("1.javascript")

	fmt.Println("2.es6")

	fmt.Println("3.es6 & typeScript")

	fmt.Println("4.es6 & typeScriptXml")

	fmt.Println("Please input 1-4")

	fmt.Scanln(&projectType)

	projectPath := tool.GetCurrentPath() + "/" + projectName

	if err := os.Mkdir(projectPath, os.ModePerm); err != nil {
		panic(err)
	}
	projectAndroid := projectPath + "/android"

	tool.ExecCmd("git", "clone", "-b", "master", "git@codeup.teambition.com:sunmi/MaxProgram/Android/Elephant.git", projectAndroid)

	if projectType == "1" {
		simpleProject(projectName,projectPath)
	}else if projectType == "2" {
		es6Project(projectName,projectPath)
	}

}
