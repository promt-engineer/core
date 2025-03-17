package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/entities"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/facade"
	"bitbucket.org/play-workspace/gocommon/tracer"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	pongWait   = 5 * time.Second
	writeWait  = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Connection struct {
	facade                   *facade.Facade
	conn                     *websocket.Conn
	srv                      *Server
	router                   *Router
	closeOnce, askOnce       *sync.Once
	responsePipeline         chan *Response
	readerClose, writerClose chan bool

	userMetaInfo *entities.PlayerMetaData

	requestMu sync.Mutex
	tr        *tracer.JaegerTracer
}

func (conn *Connection) close(wg *sync.WaitGroup) {
	conn.closeOnce.Do(func() {
		close(conn.readerClose)
		close(conn.writerClose)
		conn.conn.Close()
		wg.Done()
	})
}

func (conn *Connection) askForRemoving() {
	conn.askOnce.Do(func() {
		conn.srv.removeMe(conn)
	})
}

func (conn *Connection) reader() {
	defer func() {
		conn.askForRemoving()
	}()

	if err := conn.conn.SetReadDeadline(time.Now().Add(pongWait * 2)); err != nil {
		zap.S().Error(err)
	}

	conn.conn.SetPongHandler(conn.pongHandler)

	for {
		select {
		case _ = <-conn.readerClose:
			return
		case msg, ok := <-readMessageChan(conn.conn):
			if !ok {
				return
			}

			go conn.read(msg)
		}
	}
}

func (conn *Connection) writer() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		close(conn.responsePipeline)
		ticker.Stop()
		conn.askForRemoving()
	}()

	for {
		select {
		case resp, ok := <-conn.responsePipeline:
			conn.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				conn.conn.WriteMessage(websocket.CloseMessage, []byte{})

				return
			}

			conn.write(resp)
		case <-ticker.C:
			conn.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := conn.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case _ = <-conn.writerClose:
			return
		}
	}
}

func (conn *Connection) write(resp *Response) {
	if err := conn.conn.WriteJSON(resp); err != nil {
		zap.S().Error(err)
	}
}

func (conn *Connection) read(message []byte) {
	req := Request{}
	if err := json.Unmarshal(message, &req); err != nil {
		conn.responsePipeline <- BadRequest(err)

		return
	}

	conn.submitRequest(req)
}

func (conn *Connection) pongHandler(str string) error {
	return conn.conn.SetReadDeadline(time.Now().Add(pongWait))
}

func (conn *Connection) submitRequest(req Request) {
	conn.requestMu.Lock()
	hf, ok := conn.router.find(req.Action)
	conn.requestMu.Unlock()

	if !ok {
		conn.responsePipeline <- NotFound(req.UUID)

		return
	}

	requestBody, err := json.Marshal(req.Payload)
	if err != nil {
		conn.responsePipeline <- BadRequest(err)

		return
	}

	ctx, span := conn.tr.Start(context.Background(), "server", req.Action,
		tracer.CtxWithTraceValue|tracer.CtxWithGRPCMetadata)

	hf(HandlerBag{
		Payload:          req.Payload,
		Ctx:              ctx,
		PlayerMetaData:   conn.userMetaInfo.CopyAndSetRequest(requestBody),
		ResponsePipeline: conn.responsePipeline,
		UUID:             req.UUID})

	span.End()
}
