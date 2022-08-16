package router

import (
	"eat-and-go/handler/sd"
	"net/http"

	_ "eat-and-go/docs"
	m "eat-and-go/gorm"
	"eat-and-go/router/middleware"
	"eat-and-go/service/api"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitEngine() *gin.Engine {
	g := gin.New()
	m.DB.Init()
	return g
}

// Load loads the middlewares, routes, handler.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// The health check handler
	svcdRouter := g.Group("/sd")
	{
		svcdRouter.GET("/health", sd.HealthCheck)
		svcdRouter.GET("/version", sd.VersionCheck)
		svcdRouter.GET("/disk", sd.DiskCheck)
		svcdRouter.GET("/cpu", sd.CPUCheck)
		svcdRouter.GET("/ram", sd.RAMCheck)
	}
	main := g.Group("/api")
	//main.GET("/slack-test", api.Test)
	v1 := main.Group("/v2")
	{
		v1.POST("/slack", api.DispatchSlackEvent)
		v1.POST("/actions", api.ButtonAttachmentsHandler)
	}
	//v1.Use(middleware.SlackVerifyMiddleware())
	return g
}
