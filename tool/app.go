package tool

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"tiny_tool/module"
)

func GetAppConfig() *module.BuildConfig {
	return DeCodeAppJson(GetCurrentPath() + "/" + module.TINY_JSON)
}

func DeCodeAppJson(appJson string) *module.BuildConfig {
	if file, err := ioutil.ReadFile(appJson); err == nil {
		decoder := json.NewDecoder(bytes.NewBuffer(file))
		buildConfig := module.BuildConfig{}
		if err := decoder.Decode(&buildConfig); err == nil {
			return &buildConfig
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}

func GetApplicationId() string {
	config := GetAppConfig()
	if config == nil {
		panic("app config not found")
	} else {
		return config.Build.ApplicationId
	}

}
