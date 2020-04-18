package web_utils

import (
	"github.com/gin-gonic/gin"
	"github.com/pestanko/gouthy/app/core"
)

type HandlerFunc func(ctx *ControllerContext) error

type Controller interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type RouterWrapper struct {
	Route *gin.RouterGroup
	App   *core.GouthyApp
}

func NewRouteController(route *gin.RouterGroup, app *core.GouthyApp) RouterWrapper {
	return RouterWrapper{Route: route, App: app}
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (r *RouterWrapper) GET(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	funcs := r.wrapFunctions(handlers)

	return r.Route.GET(relativePath, funcs...)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (r *RouterWrapper) POST(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	funcs := r.wrapFunctions(handlers)

	return r.Route.POST(relativePath, funcs...)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (r *RouterWrapper) PATCH(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	funcs := r.wrapFunctions(handlers)

	return r.Route.PATCH(relativePath, funcs...)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (r *RouterWrapper) PUT(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	funcs := r.wrapFunctions(handlers)

	return r.Route.PUT(relativePath, funcs...)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (r *RouterWrapper) HEAD(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	funcs := r.wrapFunctions(handlers)

	return r.Route.HEAD(relativePath, funcs...)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (r *RouterWrapper) DELETE(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	funcs := r.wrapFunctions(handlers)

	return r.Route.DELETE(relativePath, funcs...)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (r *RouterWrapper) OPTIONS(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	funcs := r.wrapFunctions(handlers)

	return r.Route.OPTIONS(relativePath, funcs...)
}

func (r *RouterWrapper) wrapFunctions(handlers []HandlerFunc) []gin.HandlerFunc {
	var funcs []gin.HandlerFunc

	for _, f := range handlers {
		funcs = append(funcs, r.createFuncHandler(f))
	}
	return funcs
}

func (r *RouterWrapper) createFuncHandler(f HandlerFunc) func(gc *gin.Context) {
	return func(gc *gin.Context) {
		ctx, err := r.CreateControllerContext(gc)

		if err = f(ctx); err != nil {

		}
	}
}

func (r *RouterWrapper) CreateControllerContext(gc *gin.Context) (*ControllerContext, error) {
	return &ControllerContext{
		Gin: gc,
		App: r.App,
	}, nil
}
