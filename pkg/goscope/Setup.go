package goscope

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Setup is the necessary step to enable GoScope in an application.
// It will setup the necessary routes and middlewares for GoScope to work.
func Setup(config *InitData) {
	if config == nil {
		panic("Please provide a pointer to a valid and instantiated GoScopeInitData.")
	}

	ConfigSetup(config.Config)
	DatabaseSetup(DatabaseInformation{
		Type:                  Config.GoScopeDatabaseType,
		Connection:            Config.GoScopeDatabaseConnection,
		MaxOpenConnections:    Config.GoScopeDatabaseMaxOpenConnections,
		MaxIdleConnections:    Config.GoScopeDatabaseMaxIdleConnections,
		MaxConnectionLifetime: Config.GoScopeDatabaseMaxConnLifetime,
	})

	config.Router.Use(gin.Logger())
	config.Router.Use(gin.Recovery())

	logger := &LoggerGoScope{}
	gin.DefaultErrorWriter = logger

	log.SetFlags(log.Lshortfile)
	log.SetOutput(logger)

	// Use the logging middleware
	config.Router.Use(ResponseLogger)

	// Catch 404s
	config.Router.NoRoute(NoRouteResponseLogger)

	// SPA routes
	if !Config.HasFrontendDisabled {

	}

	// GoScope API
	apiGroup := config.RouteGroup.Group("/api")
	apiGroup.GET("/application-name", GetAppName)
	apiGroup.GET("/logs", getLogListHandler)
	apiGroup.GET("/requests/:id", showRequestDetailsHandler)
	apiGroup.GET("/logs/:id", showLogDetailsHandler)
	apiGroup.GET("/requests", getRequestListHandler)
	apiGroup.POST("/search/requests", searchRequestHandler)
	apiGroup.POST("/search/logs", searchLogHandler)
	apiGroup.GET("/info", getSystemInfoHandler)
}
