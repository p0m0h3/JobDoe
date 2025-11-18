package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRouteHandlers(r *gin.RouterGroup, e *Env) {
	r.GET("/ping", e.HealthzHandler)
	r.POST("/job", e.RunHandler)
	r.GET("/job/output", e.GetJobOutput)
}
