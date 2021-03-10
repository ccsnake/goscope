package main

import (
	"net/http"

	"github.com/averageflow/goscope/v3/pkg/goscope"
	"github.com/gin-gonic/gin"
)

// Setup your custom functions for the templates here
var myFunctionMap = map[string]interface{}{
	"MultiplyNumbers": func(a, b int) int { return a * b },
}

func main() {
	// Initialize an empty gin.Engine
	router := gin.New()
	// Add your custom functions to the function map
	for i := range myFunctionMap {
		router.FuncMap[i] = myFunctionMap[i]
	}
	// Setup GoScope
	applicationTemplateEngine := goscope.Setup(&goscope.InitData{
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
		},
	})

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

	// Setup any remaining routes for your application
	router.GET("/test", func(context *gin.Context) {
		context.HTML(http.StatusOK, "example.gohtml", nil)
	})

	// Start the server
	_ = router.Run(":7011")
}
