package misc

import (
	"encoding/json"
	"errors"
)

type JavaMotdJson17 struct {
	Version struct {
		Name              string `json:"name,omitempty"`
		Protocol          int    `json:"protocol,omitempty"`
		SupportedVersions []int  `json:"supportedVersions,omitempty"`
	} `json:"version"`
	Players struct {
		Max    int      `json:"max,omitempty"`
		Online int      `json:"online,omitempty"`
		Sample []Sample `json:"sample,omitempty"`
	} `json:"players"`
	Description DescriptionObjectOrString `json:"description,omitempty"`
	Modinfo     struct {
		Type    string `json:"type,omitempty"`
		ModList []struct {
			Modid   string `json:"modid,omitempty"`
			Version string `json:"version,omitempty"`
		} `json:"modList,omitempty"`
	} `json:"modinfo,omitempty"`
	Favicon             string `json:"favicon,omitempty"`
	EnforcesSecureChat  bool   `json:"enforcesSecureChat,omitempty,omitempty"`
	PreventsChatReports bool   `json:"preventsChatReports,omitempty"`
	Latency             int    `json:"latency,omitempty"`
}

type Sample struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

type DescriptionObjectOrString struct {
	DescriptionObject Description `json:"description,omitempty"`
	DescriptionString string
}
type Description struct {
	Text     string                `json:"text,omitempty"`
	ExtraObj []ExtraObjectOrString `json:"extra,omitempty"`
}

func (obj *DescriptionObjectOrString) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &obj.DescriptionString); err == nil {
		return nil
	}
	if err := json.Unmarshal(data, &obj.DescriptionObject); err == nil {
		return nil
	}
	return errors.New("invalid JSON: not a string or description object")
}

type ExtraObjectOrString struct {
	ExtraObject Extra `json:"extra,omitempty"`
	ExtraString string
}
type Extra struct {
	Text  string                `json:"text,omitempty"`
	Color string                `json:"color,omitempty"`
	Bold  bool                  `json:"bold,omitempty"`
	Extra []ExtraObjectOrString `json:"extra,omitempty"` // 递归
}

func (obj *ExtraObjectOrString) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &obj.ExtraString); err == nil {
		return nil
	}
	if err := json.Unmarshal(data, &obj.ExtraObject); err == nil {
		return nil
	}
	return errors.New("invalid JSON: not a string or description object")
}

type BedrockMotd struct {
	GameName        string `json:"gameName,omitempty"`
	HostName        string `json:"hostname,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	Version         string `json:"version,omitempty"`
	Players         int    `json:"players,omitempty"`
	MaxPlayers      int    `json:"max_players,omitempty"`
	ServerId        int64  `json:"serverId,omitempty"`
	Map             string `json:"map,omitempty"`
	GameMode        string `json:"game_mode,omitempty"`
	NintendoLimited bool   `json:"nintendo_limited,omitempty"`
	IPv4Port        int    `json:"ipv4port,omitempty"`
	IPv6Port        int    `json:"ipv6port,omitempty"`
	Extra           string `json:"extra,omitempty"`
	Latency         int    `json:"latency,omitempty"`
}

type JavaMotd16 struct {
	Protocol      string `json:"protocol,omitempty"`
	Version       string `json:"version,omitempty"`
	MessageDaily  string `json:"message_daily,omitempty"`
	OnlinePlayers int    `json:"online_players,omitempty"`
	MaxPlayers    int    `json:"max_players,omitempty"`
	Latency       int    `json:"latency,omitempty"`
}

type JavaMotd13 struct {
	Motd          string `json:"motd,omitempty"`
	OnlinePlayers int    `json:"online_players,omitempty"`
	MaxPlayers    int    `json:"max_players,omitempty"`
	Latency       int    `json:"latency,omitempty"`
}
