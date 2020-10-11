package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func NewIndexController(tools *web_utils.Tools) *IndexController {
	return &IndexController{
		Http: tools,
	}
}

type IndexController struct {
	Http *web_utils.Tools
}

func (c *IndexController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/", c.indexPage)

}

func (c *IndexController) indexPage(context *gin.Context) {
	ctx := c.Http.NewControllerContext(context)
	c.Http.HTML(ctx, http.StatusOK, "index.html", gin.H{
		"user": c.Http.GetLoggedInUser(ctx),
	})
}
