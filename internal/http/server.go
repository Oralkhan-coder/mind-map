package http

import (
	"context"
	"net/http"

	"github.com/Oralkhan-coder/mind-map/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type SimpleServer struct {
	server *gin.Engine
}

func NewSimpleServer() *SimpleServer {
	// init handlers

	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.CORP())
	router.Use(middleware.ErrorHandler())

	router.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	// protected endpoints

	return &SimpleServer{
		server: router,
	}
}

func (s *SimpleServer) Run(ctx context.Context) {
	err := s.server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
