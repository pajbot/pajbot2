package web

type user struct {
	Name                string `json:"name"`
	DisplayName         string `json:"display_name"`
	Points              int64  `json:"points"`
	Level               int64  `json:"level"`
	TotalMessageCount   int64  `json:"total_message_count"`
	OnlineMessageCount  int64  `json:"online_message_count"`
	OfflineMessageCount int64  `json:"offline_message_count"`
	LastSeen            string `json:"last_seen"`
	LastActive          string `json:"last_active"`
}

type customPayload struct {
	data map[string]interface{}
}

func (p *customPayload) Add(key string, value interface{}) {
	if p.data == nil {
		p.data = make(map[string]interface{})
	}
	p.data[key] = value
}
