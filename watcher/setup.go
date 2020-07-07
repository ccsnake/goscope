package watcher

import (
	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine) {
	// Use the logging middleware
	engine.Use(RequestLogger)
	engine.Use(ResponseLogger)

	// Setup necessary routes
	watcherGroup := engine.Group("/watcher")
	watcherGroup.GET("/", Dashboard)
	watcherGroup.GET("/requests", GetRequests)
	watcherGroup.GET("/requests/:id", ShowRequest)
	watcherGroup.GET("/responses", GetResponses)
	watcherGroup.GET("/responses/:id", ShowResponse)
}
