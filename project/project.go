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
	"tiny_tool/module"
	"tiny_tool/tool"
)

func InitProject() {
	makeProject()
}

func makeProject() {
	projectName := ""
	projectType := ""
	titlePrinter := pterm.NewStyle(pterm.FgLightCyan, pterm.BgDefault, pterm.Bold)
	titlePrinter.Println("> Please enter a project name")

	fmt.Println()

	fmt.Scanln(&projectName)

	log.Clean()

	log.Header()

	titlePrinter.Println("> Please select project type")

	pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
		{Level: 0, Text: "1.javascript", TextStyle: pterm.NewStyle(pterm.FgBlue), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 1, Text: "Basic js engineering", TextStyle: pterm.NewStyle(pterm.FgBlue), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 0, Text: "2.es6 🔥", TextStyle: pterm.NewStyle(pterm.FgGreen), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 1, Text: "Support ES6 features, it is recommended to use", TextStyle: pterm.NewStyle(pterm.FgGreen), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 0, Text: "3.jsx", TextStyle: pterm.NewStyle(pterm.FgCyan), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 1, Text: "JSX is a JavaScript syntax extension that looks a lot like XML", TextStyle: pterm.NewStyle(pterm.FgCyan), BulletStyle: pterm.NewStyle(pterm.FgRed)},
	}).Render()

	titlePrinter.Println("> Please input 1-3")

	fmt.Println()

	fmt.Scanln(&projectType)

	fmt.Println()

	log.Header()

	log.Clean()

	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("Waiting for ...")

	projectPath := tool.GetCurrentPath() + "/" + projectName

	if projectType == "1" {
		tool.BaseCmd("git", false, "clone", "-b", "v0.0.2", "--depth=1", "git@github.com:Tiny-UI/TinyJSTemplate.git", projectPath)
		introSpinner.Stop()
	} else if projectType == "2" {
		tool.BaseCmd("git", false, "clone", "-b", "v0.0.2", "--depth=1", "git@github.com:Tiny-UI/TinyES6Template.git", projectPath)
		introSpinner.Stop()
		dep.Install(projectPath + "/webpack")
	} else if projectType == "3" {
		tool.BaseCmd("git", false, "clone", "-b", "v0.0.2", "--depth=1", "git@github.com:Tiny-UI/TinyJSXTemplate.git", projectPath)
		introSpinner.Stop()
		dep.Install(projectPath + "/webpack")
	}

	introSpinner, _ = pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("Waiting for ...")

	projectAndroid := projectPath + "/android"

	tool.BaseCmd("git", false, "clone", "-b", "develop", "git@codeup.teambition.com:sunmi/Android/Tiny-UI/TinyUI.git", projectAndroid)

	writeAppConfig(projectName, projectName)

	os.RemoveAll(projectPath + "/.git")

	introSpinner.Stop()
	log.Clean()
	titlePrinter.Println("create success")
}

func writeAppConfig(projectName, projectPath string) {
	appJson := projectPath + "/" + module.TINY_JSON
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
