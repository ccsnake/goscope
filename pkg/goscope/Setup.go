package goscope

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/averageflow/goscope/v3/web"

	"github.com/averageflow/goscope/v3/internal/utils"

	"github.com/gin-gonic/gin"
)

func PrepareTemplateEngine(d *InitData) *template.Template {
	var applicationFunctionMap = map[string]interface{}{
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
		"FieldHasContent": func(fieldContent string) bool {
			return fieldContent != "" && strings.TrimSpace(fieldContent) != ""
		},
		"ResponseStatusColor": utils.ResponseStatusColor,
	}

	for i := range applicationFunctionMap {
		d.Router.FuncMap[i] = applicationFunctionMap[i]
	}

	applicationTemplateEngine := template.Must(template.New("").
		Funcs(d.Router.FuncMap).
		ParseFS(
			web.TemplateFiles,
			"templates/goscope-components/*",
			"templates/goscope-views/*",
		))

	return applicationTemplateEngine
}

// PrepareMiddleware is the necessary step to enable GoScope in an application.
// It will setup the necessary routes and middlewares for GoScope to work.
func PrepareMiddleware(d *InitData) {
	if d == nil {
		panic("Please provide a pointer to a valid and instantiated GoScopeInitData.")
	}

	configSetup(d.Config)
	databaseSetup(databaseInformation{
		databaseType:          Config.GoScopeDatabaseType,
		connection:            Config.GoScopeDatabaseConnection,
		maxOpenConnections:    Config.GoScopeDatabaseMaxOpenConnections,
		maxIdleConnections:    Config.GoScopeDatabaseMaxIdleConnections,
		maxConnectionLifetime: Config.GoScopeDatabaseMaxConnLifetime,
	})

	d.Router.Use(gin.Logger())
	d.Router.Use(gin.Recovery())

	logger := &loggerGoScope{}
	gin.DefaultErrorWriter = logger

	log.SetFlags(log.Lshortfile)
	log.SetOutput(logger)

	// Use the logging middleware
	d.Router.Use(responseLogger)

	// Catch 404s
	d.Router.NoRoute(noRouteResponseLogger)

	// SPA routes
	if !Config.HasFrontendDisabled {
		d.RouteGroup.GET("/", requestListPageHandler)
		d.RouteGroup.GET("", requestListPageHandler)
		d.RouteGroup.GET("/requests", requestListPageHandler)
		d.RouteGroup.GET("/logs", logListPageHandler)
		d.RouteGroup.GET("/logs/:id", logDetailsPageHandler)
		d.RouteGroup.GET("/requests/:id", requestDetailsPageHandler)
		d.RouteGroup.GET("/info", systemInfoPageHandler)

		d.RouteGroup.GET("/styles/:filename", func(c *gin.Context) {
			var routeData fileByRoute

			err := c.BindUri(&routeData)
			if err != nil {
				log.Println(err.Error())
				return
			}

			file, err := web.StyleFiles.ReadFile(fmt.Sprintf("styles/%s", routeData.FileName))
			if err != nil {
				log.Println(err.Error())
				return
			}

			c.Header("Content-Type", "text/css; charset=utf-8")
			c.String(http.StatusOK, string(file))
		})

		d.RouteGroup.GET("/scripts/:filename", func(c *gin.Context) {
			var routeData fileByRoute

			err := c.BindUri(&routeData)
			if err != nil {
				log.Println(err.Error())
				return
			}

			file, err := web.ScriptFiles.ReadFile(fmt.Sprintf("scripts/%s", routeData.FileName))
			if err != nil {
				log.Println(err.Error())
				return
			}

			c.Header("Content-Type", "application/javascript; charset=utf-8")
			c.String(http.StatusOK, string(file))
		})
	}

	// GoScope API
	apiGroup := d.RouteGroup.Group("/api")
	apiGroup.GET("/application-name", getAppName)
	apiGroup.GET("/logs", getLogListHandler)
	apiGroup.GET("/requests/:id", showRequestDetailsHandler)
	apiGroup.GET("/logs/:id", showLogDetailsHandler)
	apiGroup.GET("/requests", getRequestListHandler)
	apiGroup.POST("/search/requests", searchRequestHandler)
	apiGroup.POST("/search/logs", searchLogHandler)
	apiGroup.GET("/info", getSystemInfoHandler)
}
