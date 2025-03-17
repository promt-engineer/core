package middlewares

import (
	"bitbucket.org/play-workspace/gocommon/tracer"
	"github.com/gin-gonic/gin"
)

func TraceMiddleware(tr *tracer.JaegerTracer) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		trCtx, span := tr.Start(ctx.Request.Context(), "server", ctx.Request.RequestURI,
			tracer.CtxWithTraceValue|tracer.CtxWithGRPCMetadata)

		ctx.Request = ctx.Request.WithContext(trCtx)

		ctx.Next()
		span.End()
	}
}
