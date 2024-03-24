package apiServer

import "encoding/json"

type WSMessageType string

const (
	WSMessageTypeStart   WSMessageType = "start"
	WSMessageTypeStop                  = "stop"
	WSMessageTypeCommand               = "command"
	WSMessageTypeOutput
)

type WSMessage struct {
	Type    WSMessageType
	Payload json.RawMessage
}
