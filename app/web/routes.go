package web

func RegisterRoutes(s *GouthyWebServer) {
	// Templates
	s.Router.LoadHTMLGlob("resources/templates/**/*")

	// Static
	s.Router.Static("/assets", "resources/assets")
	s.Router.StaticFile("/favicon.svg", "./resources/assets/img/favicon.svg")

	// API
	apiRoute := s.Router.Group("/api")
	v1Route := apiRoute.Group("/v1")

	s.controllers.api.auth.RegisterRoutes(v1Route.Group("/auth"))
	s.controllers.api.users.RegisterRoutes(v1Route.Group("/users"))
	s.controllers.api.apps.RegisterRoutes(v1Route.Group("/applications"))

	// Well Known
	s.controllers.WellKnown.RegisterRoutes(s.Router.Group("./.well-known"))

	// pages
	pagesRoute := s.Router.Group("/pages")
	s.controllers.pages.login.RegisterRoutes(pagesRoute)
	s.controllers.pages.registration.RegisterRoutes(pagesRoute)
	s.controllers.pages.password.RegisterRoutes(pagesRoute.Group("/password"))
	s.controllers.pages.oauth2.RegisterRoutes(pagesRoute.Group("/oauth2"))
	s.controllers.pages.debug.RegisterRoutes(pagesRoute)
	s.controllers.pages.index.RegisterRoutes(s.Router.Group("/"))
}


