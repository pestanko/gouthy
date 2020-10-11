package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func NewDebugController(tools *web_utils.Tools) *DebugController {
	return &DebugController{
		Tools: tools,
	}
}

type DebugController struct {
	Tools *web_utils.Tools
}

func (c *DebugController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/_debug", c.debugPage)

}

func (c *DebugController) debugPage(context *gin.Context) {
	ctx := c.Tools.NewControllerContext(context)
	file := context.Query("file")
	c.Tools.HTML(ctx, http.StatusOK, file, gin.H{})
}
