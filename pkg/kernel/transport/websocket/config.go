package websocket

import (
	"time"
)

type Config struct {
	MaxProcessingTime time.Duration
	ReadBufferSize    int
	WriteBufferSize   int
}
