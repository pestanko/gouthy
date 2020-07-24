package pages

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/web/web_utils"
	"net/http"
)

func NewDebugController(tools *web_utils.HTTPTools) *DebugController {
	return &DebugController{
		Tools: tools,
	}
}

type DebugController struct {
	Tools *web_utils.HTTPTools
}

func (c *DebugController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/_debug", c.debugPage)

}

func (c *DebugController) debugPage(context *gin.Context) {
	ctx := c.Tools.NewControllerContext(context)
	file := context.Query("file")
	c.Tools.HTML(ctx, http.StatusOK, file, gin.H{})
}
