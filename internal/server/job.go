package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateJobRequest struct {
	Code string            `json:"code" binding:"required"`
	Env  map[string]string `json:"env"`
}

func (e Env) CreateJob(c *gin.Context) {
	req := CreateJobRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuid := e.Runner.Run(req.Code, req.Env)
	c.JSON(http.StatusOK, gin.H{"message": "code executed successfully", "uuid": uuid})
}

func (e Env) GetJobOutput(c *gin.Context) {
	data, err := e.Runner.GetOutput(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(200, "text/plain", []byte(data))
}
