package main

import (
	"github.com/averageflow/goscope/v2/src/goscopetypes"

	"github.com/averageflow/goscope/v2/src/goscope"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	goscope.Setup(&goscopetypes.GoScopeInitData{
		Router:     router,
		RouteGroup: router.Group("/goscope"),
		Config: &goscopetypes.GoScopeApplicationEnvironment{
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
