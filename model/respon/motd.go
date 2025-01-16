package respon

import (
	"satellite/model/misc"
)

type MotdResponse struct {
	Edition             string        `json:"edition,omitempty"`
	MineCraftEdition    string        `json:"mine_craft_edition,omitempty"`
	MineCraftServerType string        `json:"mine_craft_server_type,omitempty"`
	Version             string        `json:"version,omitempty"`
	OnlinePlayers       int           `json:"online_players"`
	RealPlayers         int           `json:"real_players,omitempty"`
	MaxPlayers          int           `json:"max_players,omitempty"`
	Players             []misc.Sample `json:"players,omitempty"`
	Description         string        `json:"description,omitempty"`
	Favicon             string        `json:"favicon,omitempty"`
	Protocol            string        `json:"protocol,omitempty"`
	Loader              string        `json:"loader,omitempty"`
	ModCount            int           `json:"mod_count,omitempty"`
	EnforcesSecureChat  bool          `json:"enforces_secure_chat,omitempty"`
	PreventsChatReports bool          `json:"prevents_chat_reports,omitempty"`
	ServerName          string        `json:"server_name,omitempty"`
	ServerId            int64         `json:"server_id,omitempty"`
	Map                 string        `json:"map,omitempty"`
	GameMode            string        `json:"game_mode,omitempty"`
	NintendoLimited     bool          `json:"nintendo_limited,omitempty"`
	Ipv4Port            int           `json:"ipv4_port,omitempty"`
	Ipv6Port            int           `json:"ipv6_port,omitempty"`
	Extra               string        `json:"extra,omitempty"`
}
