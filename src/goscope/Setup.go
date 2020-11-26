package goscope

import (
	"log"

	"github.com/averageflow/goscope/v2/src/goscopecontrollers"
	"github.com/averageflow/goscope/v2/src/goscopetypes"
	"github.com/averageflow/goscope/v2/src/goscopeutils"

	"github.com/gin-gonic/gin"
)

// Setup is the necessary step to enable GoScope in an application.
// It will setup the necessary routes and middlewares for GoScope to work.
func Setup(settings *goscopetypes.GoScopeInitData) {
	if settings == nil {
		panic("Please provide a pointer to a valid and instantiated GoScopeInitData.")
	}

	goscopeutils.ConfigSetup(settings.Config)
	goscopeutils.DatabaseSetup(goscopeutils.DatabaseInformation{
		Type:                  goscopeutils.Config.GoScopeDatabaseType,
		Connection:            goscopeutils.Config.GoScopeDatabaseConnection,
		MaxOpenConnections:    goscopeutils.Config.GoScopeDatabaseMaxOpenConnections,
		MaxIdleConnections:    goscopeutils.Config.GoScopeDatabaseMaxIdleConnections,
		MaxConnectionLifetime: goscopeutils.Config.GoScopeDatabaseMaxConnLifetime,
	})

	settings.Router.Use(gin.Logger())
	settings.Router.Use(gin.Recovery())

	logger := &goscopecontrollers.LoggerGoScope{}
	gin.DefaultErrorWriter = logger

	log.SetFlags(log.Lshortfile)
	log.SetOutput(logger)

	// Use the logging middleware
	settings.Router.Use(goscopecontrollers.ResponseLogger)

	// Catch 404s
	settings.Router.NoRoute(goscopecontrollers.NoRouteResponseLogger)

	// SPA routes
	if !goscopeutils.Config.HasFrontendDisabled {
		settings.RouteGroup.GET("/", goscopecontrollers.ShowDashboard)
		settings.RouteGroup.GET("/logo.svg", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/js/app.js", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/js/app.js.map", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/css/app.css", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/css/dark.css", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/css/light.css", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/css/code-blocks.css", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/css/styles.css", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/favicon.ico", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/favicon-32x32.png", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/apple-touch-icon-precomposed.png", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/apple-touch-icon.png", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/favicon-16x16.png", goscopecontrollers.GetStaticFile)
		settings.RouteGroup.GET("/logs", goscopecontrollers.ShowDashboard)
		settings.RouteGroup.GET("/logs/:uuid", goscopecontrollers.ShowDashboard)
		settings.RouteGroup.GET("/requests", goscopecontrollers.ShowDashboard)
		settings.RouteGroup.GET("/requests/:uuid", goscopecontrollers.ShowDashboard)
		settings.RouteGroup.GET("/info", goscopecontrollers.ShowDashboard)
	}

	// GoScope API
	apiGroup := settings.RouteGroup.Group("/api")
	apiGroup.GET("/application-name", goscopecontrollers.GetAppName)
	apiGroup.GET("/logs", goscopecontrollers.LogList)
	apiGroup.GET("/requests/:id", goscopecontrollers.ShowRequest)
	apiGroup.GET("/logs/:id", goscopecontrollers.ShowLog)
	apiGroup.GET("/requests", goscopecontrollers.RequestList)
	apiGroup.POST("/search/requests", goscopecontrollers.SearchRequest)
	apiGroup.OPTIONS("/search/requests", goscopecontrollers.SearchRequestOptions)
	apiGroup.POST("/search/logs", goscopecontrollers.SearchLog)
	apiGroup.OPTIONS("/search/logs", goscopecontrollers.SearchLogOptions)
	apiGroup.GET("/info", goscopecontrollers.ShowSystemInfo)
}
