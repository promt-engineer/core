package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"sync"
	"time"
)

type MutexCouple struct {
	Mu        *sync.Mutex
	UpdatedAt time.Time
}

func SessionMu(scheduler *gocron.Scheduler) func(ctx *gin.Context) {
	var (
		mapMu        sync.Mutex
		requestMuMap = map[string]MutexCouple{}
	)

	var jobFunc = func() {
		mapMu.Lock()
		defer mapMu.Unlock()

		for key, couple := range requestMuMap {
			if time.Since(couple.UpdatedAt) > time.Hour {
				delete(requestMuMap, key)
			}
		}
	}

	scheduler.Every(10).Minute().Do(jobFunc)

	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if session.Get("uuid") == nil {
			session.Set("uuid", uuid.NewString())
		}

		id, ok := session.Get("uuid").(string)
		if !ok {
			id = uuid.NewString()
			session.Set("uuid", id)
		}

		session.Save()

		mapMu.Lock()
		couple, ok := requestMuMap[id]

		if !ok {
			couple.UpdatedAt = time.Now()
			couple.Mu = &sync.Mutex{}
		}

		requestMuMap[id] = couple
		mapMu.Unlock()

		couple.Mu.Lock()
		ctx.Next()
		couple.Mu.Unlock()
	}
}
