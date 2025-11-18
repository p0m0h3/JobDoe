package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/p0m0h3/jobdoe/internal/schema"
)

func (e Env) RunHandler(c *gin.Context) {
	req := schema.CreateJobRequest{}
	if err := c.BindJSON(&req); err != nil {
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
