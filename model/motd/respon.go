package motd

type MotdJavaJson struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description struct {
		Text string `json:"text"`
	} `json:"description"`
	Modinfo struct {
		Type    string `json:"type"`
		ModList []struct {
			Modid   string `json:"modid"`
			Version string `json:"version"`
		} `json:"modList"`
	} `json:"modinfo,omitempty"`
	Favicon            string `json:"favicon,omitempty"`
	EnforcesSecureChat bool   `json:"enforcesSecureChat,omitempty"`
}
