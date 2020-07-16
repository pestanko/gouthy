package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/web/web_utils"
)

func newPages(server *WebServer) pagesController {
	return pagesController{server: server, Http: server.httpTool}
}


type pagesController struct {
	server *WebServer
	Http *web_utils.HTTPTools
}

func (ctrl *pagesController) Register(route *gin.RouterGroup) {
	route.GET("/index", ctrl.IndexPage)
	route.GET("/login", ctrl.LoginPage)
	route.GET("/register", ctrl.RegisterPage)
	route.GET("/authorization-consent", ctrl.AuthorizationConsentPage)
}

func (ctrl *pagesController) LoginPage(context *gin.Context) {
	//ctx := ctrl.Http.NewControllerContext(context)

}

func (ctrl *pagesController) RegisterPage(context *gin.Context) {

}

func (ctrl *pagesController) AuthorizationConsentPage(context *gin.Context) {

}

func (ctrl *pagesController) IndexPage(context *gin.Context) {
	
}
