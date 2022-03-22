package main

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"tiny_tool/build"
	"tiny_tool/dev"
	"tiny_tool/log"
	"tiny_tool/project"
)

func main() {
	log.Clean()
	log.Header()
	fmt.Println()
	arg := os.Args[1]
	if arg == "dev" {
		dev.Dev()
	} else if arg == "build" {
		build.Build()
	}  else if arg == "init" {
		project.InitProject()
	} else if arg == "-v" {
		fmt.Println("TINY CLI VERSION:V1.0.7")
	}else if arg == "-h" {
		pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
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
		}).Render()
	}
}
