package dep

import "MockConfig/tool"

func Install(projectPath string)  {
	tool.BaseCmd("npm", false, "i", "--prefix", projectPath)
}
