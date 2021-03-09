package main

import (
	"github.com/averageflow/goscope/v3/pkg/goscope"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	goscope.Setup(&goscope.InitData{
		Router:     router,
		RouteGroup: router.Group("/goscope"),
		Config: &goscope.Environment{
			ApplicationID:                     "",
			ApplicationName:                   "",
			ApplicationTimezone:               "Europe/Amsterdam",
			GoScopeDatabaseConnection:         "",
			GoScopeDatabaseType:               "mysql",
			GoScopeEntriesPerPage:             50,
			HasFrontendDisabled:               false,
			GoScopeDatabaseMaxOpenConnections: 10,
			GoScopeDatabaseMaxIdleConnections: 5,
			GoScopeDatabaseMaxConnLifetime:    10,
		},
	})

	_ = router.Run(":7011")
}
