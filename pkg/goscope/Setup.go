package goscope

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/averageflow/goscope/v3/web"

	"github.com/averageflow/goscope/v3/internal/utils"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

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
		d.FuncMap[i] = applicationFunctionMap[i]
	}

	applicationTemplateEngine := template.Must(template.New("").
		Funcs(d.FuncMap).
		ParseFS(
			web.TemplateFiles,
			"templates/goscope-components/*",
			"templates/goscope-views/*",
		))
	d.Router.Renderer = &TemplateRenderer{templates: applicationTemplateEngine}

	return applicationTemplateEngine
}

// PrepareMiddleware is the necessary step to enable GoScope in an application.
func PrepareMiddleware(d *InitData) (*Scope, error) {
	if d == nil || d.Config == nil {
		panic("Please provide a pointer to a valid and instantiated GoScopeInitData.")
	}

	s := &Scope{Config: d.Config}

	if err := s.setupDB(); err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}

	//d.Router.Use(middleware.Logger())
	//d.Router.Use(middleware.Recover())

	//gin.DefaultErrorWriter = logger

	// Use the logging middleware
	d.Router.Use(s.responseLogger)

	// Catch 404s
	//d.Router.NoRoute(noRouteResponseLogger)
	d.Router.HTTPErrorHandler = func(err error, context echo.Context) {
		if errors.Is(err, echo.ErrNotFound) {
			_ = s.noRouteResponseLogger(context)
		} else {
			d.Router.DefaultHTTPErrorHandler(err, context)
		}
	}

	fmt.Printf("GoScope is ready to serve requests.\n")
	fmt.Printf("GoScope is using %s as the database.\n", d.Config.GoScopeDatabaseType)
	// HasFrontend is true if the user has provided a frontend.
	fmt.Printf("GoScope is using %v as the frontend.\n", d.Config.HasFrontendDisabled)
	// SPA routes
	if !s.Config.HasFrontendDisabled {
		fmt.Printf("### SPA routes enabled\n")
		d.RouteGroup.GET("/", s.requestListPageHandler)
		d.RouteGroup.GET("", s.requestListPageHandler)
		d.RouteGroup.GET("/requests", s.requestListPageHandler)
		d.RouteGroup.GET("/logs", s.logListPageHandler)
		d.RouteGroup.GET("/logs/:id", s.logDetailsPageHandler)
		d.RouteGroup.GET("/requests/:id", s.requestDetailsPageHandler)
		d.RouteGroup.GET("/info", s.systemInfoPageHandler)

		d.RouteGroup.GET("/styles/:filename", func(c echo.Context) error {
			var routeData fileByRoute

			err := c.Bind(&routeData)
			if err != nil {
				return err
			}

			file, err := web.StyleFiles.ReadFile(fmt.Sprintf("styles/%s", routeData.FileName))
			if err != nil {
				return err
			}

			c.Response().Header().Set("Content-Type", "text/css; charset=utf-8")
			return c.String(http.StatusOK, string(file))
		})

		d.RouteGroup.GET("/scripts/:filename", func(c echo.Context) error {
			var routeData fileByRoute

			err := c.Bind(&routeData)
			if err != nil {
				c.Logger().Errorf("%s", err.Error())
				return err
			}

			file, err := web.ScriptFiles.ReadFile(fmt.Sprintf("scripts/%s", routeData.FileName))
			if err != nil {
				c.Logger().Errorf("%s", err.Error())
				return err
			}

			c.Response().Header().Set("Content-Type", "application/javascript; charset=utf-8")
			return c.String(http.StatusOK, string(file))
		})
	}

	// GoScope API
	apiGroup := d.RouteGroup.Group("/api")
	apiGroup.GET("/application-name", s.getAppName)
	apiGroup.GET("/logs", s.getLogListHandler)
	apiGroup.GET("/requests/:id", s.showRequestDetailsHandler)
	apiGroup.GET("/logs/:id", s.showLogDetailsHandler)
	apiGroup.GET("/requests", s.getRequestListHandler)
	apiGroup.POST("/search/requests", s.searchRequestHandler)
	apiGroup.POST("/search/logs", s.searchLogHandler)
	apiGroup.GET("/info", s.getSystemInfoHandler)
	return s, nil
}
