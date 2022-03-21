package project

import (
	"bytes"
	"fmt"
	"github.com/pterm/pterm"
	"io/ioutil"
	"os"
	"strings"
	"tiny_tool/dep"
	"tiny_tool/log"
	"tiny_tool/tool"
)

func InitProject() {
	makeProject()
}

func makeProject() {
	log.Clean()
	projectName := ""
	projectType := ""
	titlePrinter := pterm.NewStyle(pterm.FgLightCyan, pterm.BgGray, pterm.Bold)

	titlePrinter.Println("Please enter a project name")

	printer := pterm.NewStyle(pterm.FgLightGreen, pterm.BgWhite, pterm.Italic)

	fmt.Scanln(&projectName)

	titlePrinter.Println("Please select project type")

	printer.Println("> 1.javascript")

	printer.Println("> 2.es6")

	printer.Println("> 3.jsx")

	titlePrinter.Println("> Please input 1-3")

	fmt.Scanln(&projectType)

	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("Waiting for ...")

	projectPath := tool.GetCurrentPath() + "/" + projectName

	if projectType == "1" {
		tool.BaseCmd("git", false, "clone", "-b", "v0.0.1", "--depth=1", "git@codeup.teambition.com:sunmi/TinyTemplates/TinyJSTemplate.git", projectPath)
	} else if projectType == "2" {
		tool.BaseCmd("git", false, "clone", "-b", "v0.0.1", "--depth=1", "git@codeup.teambition.com:sunmi/TinyTemplates/TinyES6Template.git", projectPath)
		introSpinner.Stop()
		dep.Install(projectPath + "/webpack")
	} else if projectType == "3" {
		tool.BaseCmd("git", false, "clone", "-b", "v0.0.1", "--depth=1", "git@codeup.teambition.com:sunmi/TinyTemplates/TinyJSXTemplate.git", projectPath)
		introSpinner.Stop()
		dep.Install(projectPath + "/webpack")
	}

	introSpinner, _ = pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("Waiting for ...")

	projectAndroid := projectPath + "/android"

	tool.BaseCmd("git", false, "clone", "-b", "feat/tiny/template", "git@codeup.teambition.com:sunmi/MaxProgram/Android/Elephant.git", projectAndroid)

	writeAppConfig(projectName, projectName)

	os.RemoveAll(projectPath + "/.git")

	introSpinner.Stop()
	log.Clean()
	titlePrinter.Println("create success")
}

func writeAppConfig(projectName, projectPath string) {
	appJson := projectPath + "/tiny.json"
	if data, err := ioutil.ReadFile(appJson); err == nil {
		appJsonStr := string(data)
		appJsonStr = strings.ReplaceAll(appJsonStr, "$PROJECT_NAME", projectName)
		appJsonStr = strings.ReplaceAll(appJsonStr, "$PACKAGE_NAME", strings.ToLower(projectName))
		os.Remove(appJson)
		if err := ioutil.WriteFile(appJson, bytes.NewBufferString(appJsonStr).Bytes(), os.ModePerm); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
}