package module

type BuildConfig struct {
	Build   BuildModule   `json:"build"`
	Runtime RuntimeModule `json:"runtime"`

	ProjectType string `json:projectType`
	AndroidDev string `json:androidDev`
}

type KeystoreModule struct {
	KeyAlias      string `json:"keyAlias"`
	KeyPassword   string `json:"keyPassword"`
	StoreFilePath string `json:"storeFilePath"`
	StorePassword string `json:"storePassword"`
}

type RuntimeModule struct {
	BaseWidth      int          `json:"baseWidth"`
	LauncherRouter string       `json:"launcherRouter"`
	Pages          []PageModule `json:"pages"`
}

type BuildModule struct {
	AppName struct {
		Default string `json:"default"`
	} `json:"appName"`
	ApplicationId string         `json:"applicationId"`
	Dependencies  []interface{}  `json:"dependencies"`
	Keystore      KeystoreModule `json:"keystore"`
	LauncherIcon  []struct {
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
}

type PageModule struct {
	Name   string `json:"name"`
	Router string `json:"router"`
	Source string `json:"source"`
}
