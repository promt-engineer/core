package services

import (
	"container/list"
	"sync"
)

type CheatsService struct {
	cheatsList *list.List
	maxSize    int
	locker     *sync.Mutex
}

var (
	cheatSrv  *CheatsService
	cheatOnce sync.Once
)

type cheatNode struct {
	session   string
	cheatBody interface{}
}

func NewCheatsService() *CheatsService {
	cheatOnce.Do(func() {
		cheatSrv = &CheatsService{
			cheatsList: list.New(),
			maxSize:    20,
			locker:     &sync.Mutex{},
		}
	})

	return cheatSrv
}

func (s *CheatsService) Add(session string, cheatBody interface{}) {
	s.locker.Lock()
	defer s.locker.Unlock()

	s.cheatsList.PushFront(&cheatNode{session: session, cheatBody: cheatBody})

	if s.cheatsList.Len() > s.maxSize {
		s.cheatsList.Remove(s.cheatsList.Back())
	}
}

func (s *CheatsService) Get(session string) (interface{}, bool) {
	s.locker.Lock()
	defer s.locker.Unlock()

	for el := s.cheatsList.Front(); el != nil; el = el.Next() {
		val := el.Value.(*cheatNode)
		if val.session == session {
			s.cheatsList.Remove(el)

			return val.cheatBody, true
		}
	}

	return nil, false
}
