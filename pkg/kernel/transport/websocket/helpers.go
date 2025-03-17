package websocket

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func readMessageChan(conn *websocket.Conn) chan []byte {
	ch := make(chan []byte)

	go func() {
		_, message, err := conn.ReadMessage()
		if err != nil {
			zap.S().Error(err)
			close(ch)

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.S().Error(err)
			}

			return
		}
		ch <- message
	}()

	return ch
}
