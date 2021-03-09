package goscope

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/averageflow/goscope/v3/internal/utils"

	"github.com/gin-gonic/gin"
)

// Setup is the necessary step to enable GoScope in an application.
// It will setup the necessary routes and middlewares for GoScope to work.
func Setup(config *InitData) {
	if config == nil {
		panic("Please provide a pointer to a valid and instantiated GoScopeInitData.")
	}

	configSetup(config.Config, config.RouteGroup.BasePath())
	databaseSetup(databaseInformation{
		databaseType:          Config.GoScopeDatabaseType,
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

	config.Router.FuncMap["EpochToTimeAgoHappened"] = utils.EpochToTimeAgoHappened
	config.Router.FuncMap["EpochToHumanReadable"] = utils.EpochToHumanReadable

	var files []string

	templateLocation := fmt.Sprintf("%sweb/templates", "../../")

	err := filepath.Walk(templateLocation, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".gohtml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err.Error())
	}
	config.Router.LoadHTMLFiles(files...)

	// SPA routes
	if !Config.HasFrontendDisabled {
		config.RouteGroup.GET("/", requestListPageHandler)
		config.RouteGroup.GET("", requestListPageHandler)
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
