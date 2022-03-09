package dep

import "tiny_tool/tool"

func Install(projectPath string)  {
	tool.BaseCmd("npm", false, "i", "--prefix", projectPath)
}
