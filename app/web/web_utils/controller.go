package web_utils

import "github.com/gin-gonic/gin"

type Controller interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type ControllerContext struct {
	Gin *gin.Context
}
