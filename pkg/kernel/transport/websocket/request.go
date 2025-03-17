package websocket

import "github.com/google/uuid"

type Request struct {
	Action  string      `json:"action"`
	UUID    uuid.UUID   `json:"uuid"`
	Payload interface{} `json:"payload"`
}
