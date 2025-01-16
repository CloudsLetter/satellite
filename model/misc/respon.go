package misc

type MotdJson struct {
	Version struct {
		Name              string `json:"name"`
		Protocol          int    `json:"protocol"`
		SupportedVersions []int  `json:"supportedVersions,omitempty"`
	} `json:"version"`
	Players struct {
		Max    int      `json:"max"`
		Online int      `json:"online"`
		Sample []Sample `json:"sample,omitempty"`
	} `json:"players"`
	Description struct {
		Text  string `json:"text"`
		Extra []struct {
			Text  string `json:"text"`
			Color string `json:"color"`
			Bold  bool   `json:"bold,omitempty"`
		} `json:"extra"`
	} `json:"description"`
	Modinfo struct {
		Type    string `json:"type"`
		ModList []struct {
			Modid   string `json:"modid"`
			Version string `json:"version"`
		} `json:"modList,omitempty"`
	} `json:"modinfo,omitempty"`
	Favicon             string `json:"favicon,omitempty"`
	EnforcesSecureChat  bool   `json:"enforcesSecureChat,omitempty,omitempty"`
	PreventsChatReports bool   `json:"preventsChatReports,omitempty"`
}

type Sample struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}
