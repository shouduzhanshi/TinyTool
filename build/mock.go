package build

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"tiny_tool/module"
	"tiny_tool/server"
	"tiny_tool/tool"
)

func CreateAndroidBuildConfig(dslDir string, out *string) string {

	appConfig := tool.GetAppConfig()

	pages := list.New()

	processPages(dslDir, appConfig, pages)

	modules := make([]module.PageModule, 0)

	for i := pages.Front(); i != nil; i = i.Next() {
		modules = append(modules, i.Value.(module.PageModule))
	}
	appConfig.Runtime.Pages = modules
	appConfig.Build.Keystore.StoreFilePath = tool.GetAbsPath(projectPath, appConfig.Build.Keystore.StoreFilePath)
	icon := appConfig.Build.LauncherIcon
	for i := 0; i < len(icon); i++ {
		icon[i].Icon = tool.GetAbsPath(projectPath, icon[i].Icon)
	}
	splash := appConfig.Build.Splash.Background
	for i := 0; i < len(splash); i++ {
		splash[i].Src = tool.GetAbsPath(projectPath, splash[i].Src)
	}
	appJsonPath := projectPath + "/.mock.json"
	appConfig.Runtime.Ws = server.GetWsPath()
	if out!=nil {
		appConfig.Build.Configs.OutApk = *out
	}
	if data, err := json.Marshal(appConfig); err != nil {
		panic(err)
	} else {
		if err := ioutil.WriteFile(appJsonPath, data, os.ModePerm); err != nil {
			panic(err)
		}
	}
	return appJsonPath
}

func processPages(dslDir string, appConfig *module.BuildConfig, pages *list.List) {
	if files, err := ioutil.ReadDir(dslDir); err == nil {
		for _, file := range files {
			fileName := file.Name()
			if file.IsDir() {
				processPages(dslDir+"/"+fileName, appConfig, pages)
			} else {
				if strings.HasSuffix(fileName, ".js") {
					modules := appConfig.Runtime.Pages
					for _, page := range modules {
						if page.Name == fileName[0:len(fileName)-3] || dslDir+"/"+fileName == tool.GetAbsPath(tool.GetCurrentPath(), page.Source) {
							page.Source = dslDir + "/" + fileName
							pages.PushBack(page)
						}
					}
				}
			}
		}
	}
}
