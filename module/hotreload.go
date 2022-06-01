package module

type HotReloadModule struct {
	Router string `json:"router"`
	Name   string `json:"name"`
	Data   string `json:"data"`
	Size   int64    `json:"size"`
	FileName string `json:"fileName"`
}
