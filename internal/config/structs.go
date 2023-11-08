package config

type T struct {
	Discord struct {
		Token  string `json:"token"`
		Status struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"status"`
	} `json:"discord"`
	Guild struct {
		Id                 string `json:"id"`
		GeneralChannel     string `json:"general-channel"`
		TicketsCategory    string `json:"tickets-category"`
		TicketsLogsChannel string `json:"tickets-logs-channel"`
		MirrorWebhook      string `json:"mirror-webhook"`
	} `json:"guild"`
	Roles struct {
		PurgeRevive string   `json:"purge-revive"`
		Verified    string   `json:"verified"`
		Revolution  string   `json:"revolution"`
		Luminary    string   `json:"luminary"`
		Staff       []string `json:"staff"`
	} `json:"roles"`
}
