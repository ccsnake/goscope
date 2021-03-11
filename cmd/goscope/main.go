package main

import (
	"net/http"

	"github.com/averageflow/goscope/v3/pkg/goscope"
	"github.com/gin-gonic/gin"
)

// PrepareMiddleware your custom functions for the templates here.
var myFunctionMap = map[string]interface{}{
	"MultiplyNumbers": func(a, b int) int { return a * b },
}

func main() {
	// Initialize an empty gin.Engine
	router := gin.New()

	goScopeConfig := goscope.InitData{
		Router:     router,
		RouteGroup: router.Group("/goscope"),
		Config: &goscope.Environment{
			ApplicationID:                     "go-scope",
			ApplicationName:                   "go-scope",
			ApplicationTimezone:               "Europe/Amsterdam",
			GoScopeDatabaseConnection:         "root:root@tcp(127.0.0.1:3306)/go_scope",
			GoScopeDatabaseType:               "mysql",
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

	// Optionally add your custom functions to the function map of the router
	for i := range myFunctionMap {
		router.FuncMap[i] = myFunctionMap[i]
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
		router.SetHTMLTemplate(applicationTemplateEngine)
	}

	// Required call to the PrepareMiddleware of GoScope
	goscope.PrepareMiddleware(&goScopeConfig)

	// PrepareMiddleware any remaining routes for your application
	router.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "example.gohtml", nil)
	})

	// Start the server
	_ = router.Run(":7011")
}
