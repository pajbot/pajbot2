package ws

import (
	"encoding/json"
	"log"
)

// Payload xD
type Payload struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

// ToJSON creates a json string from the payload
func (p *Payload) ToJSON() (ret []byte) {
	ret, err := json.Marshal(p)
	if err != nil {
		log.Println("Error marshalling payload:", err)
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
