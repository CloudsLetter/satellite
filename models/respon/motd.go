package respon

import (
	"satellite/models/minecraft"
)

type MotdResponse struct {
	Edition             string              `json:"edition,omitempty"`
	MineCraftEdition    string              `json:"mine_craft_edition,omitempty"`
	MineCraftServerType string              `json:"mine_craft_server_type,omitempty"`
	Version             string              `json:"version,omitempty"`
	OnlinePlayers       int                 `json:"online_players"`
	RealPlayers         int                 `json:"real_players,omitempty"`
	MaxPlayers          int                 `json:"max_players,omitempty"`
	Players             []minecraft.Sample  `json:"players,omitempty"`
	Description         string              `json:"description,omitempty"`
	DescriptionRaw      []DescriptionRaw    `json:"description_raw,omitempty"`
	Favicon             string              `json:"favicon,omitempty"`
	Protocol            string              `json:"protocol,omitempty"`
	Loader              string              `json:"loader,omitempty"`
	ModCount            int                 `json:"mod_count,omitempty"`
	ModList             []minecraft.ModLisS `json:"mod_list,omitempty"`
	EnforcesSecureChat  bool                `json:"enforces_secure_chat,omitempty"`
	PreventsChatReports bool                `json:"prevents_chat_reports,omitempty"`
	ServerName          string              `json:"server_name,omitempty"`
	ServerId            int64               `json:"server_id,omitempty"`
	Map                 string              `json:"map,omitempty"`
	GameMode            string              `json:"game_mode,omitempty"`
	NintendoLimited     bool                `json:"nintendo_limited,omitempty"`
	Ipv4Port            int                 `json:"ipv4_port,omitempty"`
	Ipv6Port            int                 `json:"ipv6_port,omitempty"`
	Latency             int                 `json:"latency,omitempty"`
	Extra               string              `json:"extra,omitempty"`
}
type DescriptionRaw struct {
	Text  string `json:"text,omitempty"`
	Color string `json:"color,omitempty"`
	Bold  bool   `json:"bold,omitempty"`
	Depth int    `json:"depth,omitempty"`
}
