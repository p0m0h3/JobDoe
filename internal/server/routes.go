package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRouteHandlers(r *gin.RouterGroup, e *Env) {
	r.GET("/healthz", e.HealthzHandler)
	r.POST("/job", e.CreateJob)
	r.GET("/job/output", e.GetJobOutput)
	r.POST("/volume", e.CreateVolume)
}
