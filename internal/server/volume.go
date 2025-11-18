package server

import "github.com/gin-gonic/gin"

type CreateVolumeRequest struct {
	Name string `json:"name" binding:"required"`
}

func (e Env) CreateVolume(c *gin.Context) {
	req := CreateVolumeRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := e.Volume.CreateVolume(req.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Volume created"})
}
