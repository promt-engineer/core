package websocket

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/entities"
	"context"
	"github.com/google/uuid"
)

type Handler interface {
	Register(r *Router)
	Shutdown()
}

type HandlerBag struct {
	Payload          interface{}
	Ctx              context.Context
	UUID             uuid.UUID
	PlayerMetaData   *entities.PlayerMetaData
	ResponsePipeline chan *Response
}

type HandleFunc func(HandlerBag)

type Router struct {
	routeMap map[string]HandleFunc
}

func NewRouter() *Router {
	return &Router{routeMap: map[string]HandleFunc{}}
}

func (r *Router) Accept(action string, hf HandleFunc) {
	r.routeMap[action] = hf
}

func (r *Router) find(action string) (hf HandleFunc, ok bool) {
	hf, ok = r.routeMap[action]

	return hf, ok
}
