package main

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"tiny_tool/build"
	"tiny_tool/dep"
	"tiny_tool/dev"
	"tiny_tool/log"
	"tiny_tool/project"
)

func main() {
	log.Clean()
	arg := os.Args[1]
	if arg == "dev" {
		dev.Dev()
	} else if arg == "build" {
		build.Build()
	} else if arg == "init" {
		project.InitProject()
	} else if arg == "-v" {
		fmt.Println("TINY CLI VERSION:V1.0.8")
	} else if arg == "-h" {
		showHelp()
	} else if arg == "clean" {
		project.Clean()
	} else if arg == "dep" {
		dep.Install()
	}
}

func showHelp() error {
	return pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
		{Level: 0, Text: "Further help:", TextStyle: pterm.NewStyle(pterm.FgLightBlue), BulletStyle: pterm.NewStyle(pterm.FgRed)},

		{Level: 1, Text: "dev", TextStyle: pterm.NewStyle(pterm.FgBlue), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 3, Text: "Running debug artifacts, and hot reloading", TextStyle: pterm.NewStyle(pterm.FgBlue), BulletStyle: pterm.NewStyle(pterm.FgRed)},

		{Level: 1, Text: "build", TextStyle: pterm.NewStyle(pterm.FgGreen), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 3, Text: "Build a release product", TextStyle: pterm.NewStyle(pterm.FgGreen), BulletStyle: pterm.NewStyle(pterm.FgRed)},

		{Level: 1, Text: "init", TextStyle: pterm.NewStyle(pterm.FgCyan), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 3, Text: "Create new project", TextStyle: pterm.NewStyle(pterm.FgCyan), BulletStyle: pterm.NewStyle(pterm.FgRed)},

		{Level: 1, Text: "-v", TextStyle: pterm.NewStyle(pterm.FgLightYellow), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 3, Text: "Current tool version", TextStyle: pterm.NewStyle(pterm.FgLightYellow), BulletStyle: pterm.NewStyle(pterm.FgRed)},

		{Level: 1, Text: "-h", TextStyle: pterm.NewStyle(pterm.FgLightMagenta), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 3, Text: "Show help", TextStyle: pterm.NewStyle(pterm.FgLightMagenta), BulletStyle: pterm.NewStyle(pterm.FgRed)},

		{Level: 1, Text: "clean", TextStyle: pterm.NewStyle(pterm.FgRed), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 3, Text: "Clean build cache", TextStyle: pterm.NewStyle(pterm.FgRed), BulletStyle: pterm.NewStyle(pterm.FgRed)},

		{Level: 1, Text: "dep", TextStyle: pterm.NewStyle(pterm.FgDefault), BulletStyle: pterm.NewStyle(pterm.FgRed)},
		{Level: 3, Text: "install dsl modules", TextStyle: pterm.NewStyle(pterm.FgDefault), BulletStyle: pterm.NewStyle(pterm.FgRed)},
	}).Render()
}
