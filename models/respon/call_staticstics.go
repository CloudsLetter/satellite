package respon

type CallStatistics struct {
	MinuteStruct []struct {
		Hour   int   `json:"hour"`
		Minute int   `json:"minute"`
		Count  int64 `json:"count"`
	} `json:"minute_struct,omitempty"`
	HourStruct []struct {
		Day   int   `json:"day"`
		Hour  int   `json:"hour"`
		Count int64 `json:"count"`
	} `json:"hour_struct,omitempty"`
	DayStruct []struct {
		Month int   `json:"month"`
		Day   int   `json:"day"`
		Count int64 `json:"count"`
	} `json:"day_struct,omitempty"`
	MonthStruct []struct {
		Year  int   `json:"year"`
		Month int   `json:"month"`
		Count int64 `json:"count"`
	} `json:"month_struct,omitempty"`
}
