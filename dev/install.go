package dev

import (
	"io/ioutil"
	"time"
	"tiny_tool/build"
	"tiny_tool/log"
	"tiny_tool/tool"
	"tiny_tool/tool/xml"
)

func InstallApk(success func(), fail func()) {

	list := tool.GetDeviceList()

	androidDir := build.GetAndroidDir()

	isSuccess := false
	if list != nil || len(list) > 0 {
		for _, device := range list {
			if device.Online {
				log.V("install app to ", device.Id, " ....")
				installStart := time.Now().UnixNano()
				tool.Adb("-s", device.Id, "install", "-r", androidDir+"/build/outputs/apk/debug/app-debug.apk")
				log.E("install app to ", device.Id, " duration ", (time.Now().UnixNano()-installStart)/1e6, " ms")
				openStart := time.Now().UnixNano()
				AndroidManifestPath := androidDir + "/build/intermediates/merged_manifest/debug/AndroidManifest.xml"
				if AndroidManifestData, err := ioutil.ReadFile(AndroidManifestPath); err != nil {
					panic(err)
				} else {
					splash := getSplashActivity(AndroidManifestData)
					tool.Adb("-s", device.Id, "shell", "am", "start", "-n", tool.GetApplicationId()+"/"+splash)
					log.E("open app from ", device.Id, " ", (time.Now().UnixNano()-openStart)/1e6, " ms ")
					isSuccess = true
				}
			}
		}
	}
	if isSuccess {
		if success != nil {
			success()
		}
	} else {
		if fail != nil {
			fail()
		}
	}
}

func getSplashActivity(AndroidManifestData []byte) string {
	element := getActivityElement(AndroidManifestData)
	for _, attr := range element.Attr {
		if attr.Space == "android" {
			return attr.Value
		}
	}
	panic("splash not found")
}

func getActivityElement(AndroidManifestData []byte) *xml.Element {
	doc := xml.NewDocument()
	if err := doc.ReadFromBytes(AndroidManifestData); err == nil {
		elements := doc.SelectElement("manifest").SelectElement("application").SelectElements("activity")
		for _, element := range elements {
			for _, actions := range element.SelectElements("intent-filter") {
				//element.SelectElements("intent-filter")
				actions := actions.SelectElements("action")
				for _, action := range actions {
					attr := action.Attr
					for _, attrValue := range attr {
						if "android.intent.action.MAIN" == attrValue.Value {
							return element
						}
					}
				}
			}
		}
	} else {
		panic(err)
	}
	return nil
}

//
//func ShowApkDownloadQrcode()  {
//	log.E("total duration ", time.Now().Unix()-start, " s")
//	if tool.DeviceOnline() == nil {
//		log.E("device not found!")
//		go openBrowser(appConfig, start)
//	}
//}
//
//
//func openBrowser(appConfig module.BuildConfig) {
//	for {
//		time.Sleep(time.Duration(500) * time.Millisecond)
//		if resp, err := http.Get("http://127.0.0.1:1323/qrCode"); err == nil {
//			if resp.StatusCode == 200 {
//				data := make(map[string]interface{})
//				data["type"] = "apk"
//				data["url"] = server.GetApkDownloadUrl()
//				server.PublishMsg(data, 0)
//				if !appConfig.DisableOpenBrowser {
//					tool.ExecCmd("open", "-a", "Google Chrome", "http://127.0.0.1:1323")
//				}
//				resp.Body.Close()
//				return
//			}
//			resp.Body.Close()
//		}
//	}
//}
