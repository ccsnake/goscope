package goscope

import (
	"html/template"
	"log"

	"github.com/averageflow/goscope/v3/web"

	"github.com/averageflow/goscope/v3/internal/utils"

	"github.com/gin-gonic/gin"
)

var GoScopeFunctionMap = map[string]interface{}{
	"EpochToTimeAgoHappened": utils.EpochToTimeAgoHappened,
	"EpochToHumanReadable":   utils.EpochToHumanReadable,
	"Add":                    func(a, b int) int { return a + b },
	"SubtractTillZero": func(a, b int) int {
		result := a - b
		if result < 0 {
			return a
		}

		return result
	},
}

// Setup is the necessary step to enable GoScope in an application.
// It will setup the necessary routes and middlewares for GoScope to work.
func Setup(config *InitData) *template.Template {
	if config == nil {
		panic("Please provide a pointer to a valid and instantiated GoScopeInitData.")
	}

	configSetup(config.Config, config.RouteGroup.BasePath())
	databaseSetup(databaseInformation{
		databaseType:          Config.GoScopeDatabaseType,
		connection:            Config.GoScopeDatabaseConnection,
		maxOpenConnections:    Config.GoScopeDatabaseMaxOpenConnections,
		maxIdleConnections:    Config.GoScopeDatabaseMaxIdleConnections,
		maxConnectionLifetime: Config.GoScopeDatabaseMaxConnLifetime,
	})

	config.Router.Use(gin.Logger())
	config.Router.Use(gin.Recovery())

	logger := &loggerGoScope{}
	gin.DefaultErrorWriter = logger

	log.SetFlags(log.Lshortfile)
	log.SetOutput(logger)

	// Use the logging middleware
	config.Router.Use(responseLogger)

	// Catch 404s
	config.Router.NoRoute(noRouteResponseLogger)

	for i := range GoScopeFunctionMap {
		config.Router.FuncMap[i] = GoScopeFunctionMap[i]
	}

	// SPA routes
	if !Config.HasFrontendDisabled {
		config.RouteGroup.GET("/", requestListPageHandler)
		config.RouteGroup.GET("", requestListPageHandler)
		config.RouteGroup.GET("/logs", logListPageHandler)
		config.RouteGroup.GET("/logs/:id", logDetailsPageHandler)
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

	templateEngineNew := template.Must(template.New("").
		Funcs(config.Router.FuncMap).
		ParseFS(web.TemplateFiles, "templates/goscope-components/*", "templates/goscope-views/*"))

	return templateEngineNew
}
