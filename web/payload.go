package web

import "encoding/json"

// Payload xD
type Payload struct {
	Event string            `json:"event"`
	Data  map[string]string `json:"data"`
}

// ToJSON creates a json string from the payload
func (p *Payload) ToJSON() (ret []byte) {
	ret, err := json.Marshal(p)
	if err != nil {
		log.Error("Erro marshalling payload:", err)
	}
	return
}

func createPayload(data []byte) (*Payload, error) {
	p := &Payload{}
	err := json.Unmarshal(data, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
