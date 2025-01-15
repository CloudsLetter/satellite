package respon

type MotdResponse struct {
	Name          string `json:"name"`
	OnlinePlayers int    `json:"online_players"`
	MaxPlayers    int    `json:"max_players"`
	Sample        []struct {
		Name string `json:"name,omitempty"`
		Id   string `json:"id,omitempty"`
	} `json:"sample,omitempty"`
	Description        string `json:"description,omitempty"`
	Favicon            string `json:"favicon,omitempty"`
	Protocol           int    `json:"protocol,omitempty"`
	Loader             string `json:"loader,omitempty"`
	ModCount           int    `json:"mod_count,omitempty"`
	EnforcesSecureChat bool   `json:"enforces_secure_chat,omitempty"`
}
