package module

type BuildConfig struct {
	Build struct {
		AppName struct {
			Default string `json:"default"`
		} `json:"appName"`
		ApplicationId string        `json:"applicationId"`
		Dependencies  []interface{} `json:"dependencies"`
		Keystore      struct {
			KeyAlias      string `json:"keyAlias"`
			KeyPassword   string `json:"keyPassword"`
			StoreFilePath string `json:"storeFilePath"`
			StorePassword string `json:"storePassword"`
		} `json:"keystore"`
		LauncherIcon []struct {
			Icon       string `json:"icon"`
			Resolution string `json:"resolution"`
		} `json:"launcherIcon"`
		Splash struct {
			Background []struct {
				Resolution string `json:"resolution"`
				Src        string `json:"src"`
			} `json:"background"`
		} `json:"splash"`
		VersionCode int    `json:"versionCode"`
		VersionName string `json:"versionName"`
	} `json:"build"`
	Runtime struct {
		BaseWidth      int    `json:"baseWidth"`
		LauncherRouter string `json:"launcherRouter"`
		Pages          []struct {
			Name   string `json:"name"`
			Router string `json:"router"`
			Source string `json:"source"`
		} `json:"pages"`
	} `json:"runtime"`
}
