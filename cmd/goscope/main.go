package main

import (
	"net/http"

	"github.com/averageflow/goscope/v3/pkg/goscope"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	templateEngineNew := goscope.Setup(&goscope.InitData{
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

	if templateEngineNew != nil {
		_, err := templateEngineNew.ParseFiles("../../web/example.gohtml")
		if err != nil {
			panic(err.Error())
		}

		router.SetHTMLTemplate(templateEngineNew)
	}

	router.GET("/test", func(context *gin.Context) {
		context.HTML(http.StatusOK, "example.gohtml", nil)
	})
	_ = router.Run(":7011")
}
