package main

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/averageflow/goscope/v3/pkg/goscope"
)

func main() {
	// Initialize an empty gin.Engine
	router := echo.New()

	goScopeConfig := goscope.InitData{
		Router:     router,
		RouteGroup: router.Group("/goscope"),
		FuncMap:    map[string]interface{}{},
		Config: &goscope.Environment{
			ApplicationID:                     "go-scope",
			ApplicationName:                   "go-scope",
			ApplicationTimezone:               "asia/Shanghai",
			GoScopeDatabaseConnection:         "./goscope.sqlite",
			GoScopeDatabaseType:               "sqlite3",
			GoScopeEntriesPerPage:             50,
			HasFrontendDisabled:               false,
			GoScopeDatabaseMaxOpenConnections: 10,
			GoScopeDatabaseMaxIdleConnections: 5,
			GoScopeDatabaseMaxConnLifetime:    10,
			BaseURL:                           "/goscope",
		},
	}

	// Optionally enable GoScope frontend + your own templates
	// If you skip this code of block you should also set the
	// HasFrontendDisabled option to false to avoid panics when attempting
	// to load routes that render the HTML templates.

	// add your custom functions for the templates here.
	var myFunctionMap = map[string]interface{}{
		"MultiplyNumbers": func(a, b int) int { return a * b },
	}

	// Optionally add your custom functions to the function map of the router
	for i := range myFunctionMap {
		goScopeConfig.FuncMap[i] = myFunctionMap[i]
	}

	applicationTemplateEngine := goscope.PrepareTemplateEngine(&goScopeConfig)
	// If the application template engine is valid use it for the router
	if applicationTemplateEngine != nil {
		// Parse any other template files here
		_, err := applicationTemplateEngine.ParseFiles("../../web/example.gohtml")
		if err != nil {
			panic(err.Error())
		}
		// Finally set the html renderer of the application to the GoScope + your templates engine
		goScopeConfig.SetHTMLTemplate(applicationTemplateEngine)
	}

	// Required call to the PrepareMiddleware of GoScope
	s, err := goscope.PrepareMiddleware(&goScopeConfig)
	if err != nil {
		panic(err.Error())
	}

	defer s.Close()

	// PrepareMiddleware any remaining routes for your application
	router.GET("/test", func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, "example.gohtml", nil)
	})

	// Start the server
	_ = router.Start(":7011")
}
