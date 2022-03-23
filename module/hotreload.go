package module

type HotReloadModule struct {
	Router  string `json:"router"`
	Name    string `json:"name"`
	Data    string `json:"data"`
	Editing bool   `json:"editing"`
}
