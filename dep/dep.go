package dep

import (
	"github.com/pterm/pterm"
	"tiny_tool/tool"
)

func Install()  {
	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(false).WithRemoveWhenDone(true).Start("install npm dep ...")
	tool.BaseCmd("npm", false,"i", "--prefix", tool.GetCurrentPath()+"/webpack")
	introSpinner.Stop()
}
