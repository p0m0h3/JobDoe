package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRouteHandlers(r *gin.RouterGroup, ctx HandlerContext) {
	r.GET("/ping", ctx.HealthzHandler)
	r.POST("/job", ctx.RunHandler)
	r.GET("/job/output", ctx.GetJobOutput)
}
