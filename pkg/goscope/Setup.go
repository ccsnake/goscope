package goscope

import (
	"log"

	"github.com/averageflow/goscope/v3/internal/controllers"

	"github.com/gin-gonic/gin"
)

// Setup is the necessary step to enable GoScope in an application.
// It will setup the necessary routes and middlewares for GoScope to work.
func Setup(settings *InitData) {
	if settings == nil {
		panic("Please provide a pointer to a valid and instantiated GoScopeInitData.")
	}

	ConfigSetup(settings.Config)
	DatabaseSetup(DatabaseInformation{
		Type:                  Config.GoScopeDatabaseType,
		Connection:            Config.GoScopeDatabaseConnection,
		MaxOpenConnections:    Config.GoScopeDatabaseMaxOpenConnections,
		MaxIdleConnections:    Config.GoScopeDatabaseMaxIdleConnections,
		MaxConnectionLifetime: Config.GoScopeDatabaseMaxConnLifetime,
	})

	settings.Router.Use(gin.Logger())
	settings.Router.Use(gin.Recovery())

	logger := &controllers.LoggerGoScope{}
	gin.DefaultErrorWriter = logger

	log.SetFlags(log.Lshortfile)
	log.SetOutput(logger)

	// Use the logging middleware
	settings.Router.Use(controllers.ResponseLogger)

	// Catch 404s
	settings.Router.NoRoute(controllers.NoRouteResponseLogger)

	// SPA routes
	if !Config.HasFrontendDisabled {

	}

	// GoScope API
	apiGroup := settings.RouteGroup.Group("/api")
	apiGroup.GET("/application-name", controllers.GetAppName)
	apiGroup.GET("/logs", controllers.LogList)
	apiGroup.GET("/requests/:id", controllers.ShowRequest)
	apiGroup.GET("/logs/:id", controllers.ShowLog)
	apiGroup.GET("/requests", controllers.RequestList)
	apiGroup.POST("/search/requests", controllers.SearchRequest)
	apiGroup.OPTIONS("/search/requests", controllers.SearchRequestOptions)
	apiGroup.POST("/search/logs", controllers.SearchLog)
	apiGroup.OPTIONS("/search/logs", controllers.SearchLogOptions)
	apiGroup.GET("/info", controllers.ShowSystemInfo)
}
