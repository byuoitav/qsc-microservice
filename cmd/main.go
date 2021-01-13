package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var port int
	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	dsps := &sync.Map{}
	cfg := zap.NewProductionConfig{}
	cfg.Level.SetLevel(zapcore.DebugLevel)

	handlers := Handlers{
		CreateDSP: func(addr string) *qsc.DSP {
			if dsp, ok := dsps.Load(addr); ok {
				return dsp.(*qsc.DSP)
			}

			dsp := qsc.New(addr)

			dsps.Store(addr, dsp)
			return dsp
		},
	}

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	api := e.Group("/api/v1")
	handlers.RegisterRoutes(api)

	log.Printf("Server started on %v", lis.Addr())
	if err := e.Server.Serve(lis); err != nil {
		log.Printf("unable to serve: %s", err)
	}
}
