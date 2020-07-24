package pages

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/web/web_utils"
	"net/http"
)

func NewIndexController(tools *web_utils.HTTPTools) *IndexController {
	return &IndexController{
		Tools: tools,
	}
}

type IndexController struct {
	Tools *web_utils.HTTPTools
}

func (c *IndexController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/", c.indexPage)

}

func (c *IndexController) indexPage(context *gin.Context) {
	ctx := c.Tools.NewControllerContext(context)
	c.Tools.HTML(ctx, http.StatusOK, "index.html", gin.H{})
}
