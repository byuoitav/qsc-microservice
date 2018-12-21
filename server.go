package main

import (
	"net/http"

	"github.com/byuoitav/authmiddleware"
	"github.com/byuoitav/qsc-microservice/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := ":8016"
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	//functionality
	secure.GET("/:address/:name/volume/mute", handlers.Mute)
	secure.GET("/:address/:name/volume/unmute", handlers.UnMute)
	secure.GET("/:address/:name/volume/set/:level", handlers.SetVolume)

	//status
	secure.GET("/:address/:name/volume/level", handlers.GetVolume)
	secure.GET("/:address/:name/mute/status", handlers.GetMute)

	secure.PUT("/:address/generic/:name/:value", handlers.SetGeneric)
	secure.GET("/:address/generic/:name", handlers.GetGeneric)
	secure.GET("/:address/hardware", handlers.GetInfo)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
