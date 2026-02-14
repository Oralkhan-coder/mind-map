package http

import (
	"context"

	"github.com/Oralkhan-coder/mind-map/config"
	"github.com/Oralkhan-coder/mind-map/internal/http/middleware"
	"github.com/Oralkhan-coder/mind-map/internal/http/transport"
	"github.com/gin-gonic/gin"
)

type SimpleServer struct {
	server  *gin.Engine
	authSrv AuthSrv
}

func NewSimpleServer(auth AuthSrv, cfg *config.SecretConfig) *SimpleServer {
	authHandler := transport.NewAuthHandler(auth)

	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.CORP())
	router.Use(middleware.ErrorHandler())

	router.POST("/signup", authHandler.SignUp)
	router.POST("/login", authHandler.Login)
	router.GET("/confirm", authHandler.ConfirmEmail)
	router.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"hello": "world"}) })

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JwtSecret))
	{
		protected.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"ping": "pong"}) })
	}

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
