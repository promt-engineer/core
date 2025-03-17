package websocket

import (
	"net/http"
	"sync"

	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/entities"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/facade"
	"bitbucket.org/play-workspace/gocommon/tracer"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Server struct {
	conf     *Config
	upgrader websocket.Upgrader
	pool     map[*Connection]bool
	facade   *facade.Facade
	router   *Router
	handlers []Handler

	wg *sync.WaitGroup
	tr *tracer.JaegerTracer
}

func NewServer(conf *Config, facade *facade.Facade, tr *tracer.JaegerTracer, handlers ...Handler) *Server {
	srv := &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  conf.ReadBufferSize,
			WriteBufferSize: conf.WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		conf:   conf,
		pool:   map[*Connection]bool{},
		facade: facade,

		router:   NewRouter(),
		handlers: handlers,
		wg:       &sync.WaitGroup{},
		tr:       tr,
	}

	for _, handler := range handlers {
		handler.Register(srv.router)
	}

	return srv
}

func (s *Server) Shutdown() {
	for conn := range s.pool {
		s.removeMe(conn)
	}

	for _, handler := range s.handlers {
		handler.Shutdown()
	}

	s.wg.Wait()
	zap.S().Info("Websocket server is shutdown")
}

func (s *Server) ServeWS(ctx *gin.Context) error {
	md := entities.NewPlayerMetaDataFromCtx(ctx, nil)

	conn, err := s.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return err
	}

	realConn := &Connection{
		conn:   conn,
		facade: s.facade,

		responsePipeline: make(chan *Response),
		router:           s.router,

		readerClose: make(chan bool),
		writerClose: make(chan bool),

		askOnce:   &sync.Once{},
		closeOnce: &sync.Once{},

		userMetaInfo: md,

		srv: s,
		tr:  s.tr,
	}

	s.pool[realConn] = true
	s.wg.Add(1)

	go realConn.reader()
	go realConn.writer()

	return nil
}

func (s *Server) removeMe(conn *Connection) {
	delete(s.pool, conn)
	conn.close(s.wg)
}
