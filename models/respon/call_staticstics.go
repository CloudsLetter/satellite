package respon

type CallStatistics struct {
	Day        int   `json:"day"`
	Month      int   `json:"month"`
	Calls      int64 `json:"calls"`
	HourStruct []struct {
		Hour  int   `json:"hour"`
		Calls int64 `json:"calls"`
	} `json:"hour_struct,omitempty"`
}
