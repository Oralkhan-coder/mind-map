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
	mapSrv  MapSrv
}

func NewSimpleServer(auth AuthSrv, maps MapSrv, cfg *config.SecretConfig) *SimpleServer {
	authHandler := transport.NewAuthHandler(auth)
	mapHandler := transport.NewMapHandler(maps)

	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.CORP())
	router.Use(middleware.ErrorHandler())

	router.POST("/signup", authHandler.SignUp)
	router.POST("/login", authHandler.Login)
	router.GET("/confirm", authHandler.ConfirmEmail)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JwtSecret))
	{
		protected.GET("/maps", mapHandler.GetMaps)
		protected.POST("/maps", mapHandler.CreateMap)
		protected.GET("/maps/:id", mapHandler.GetByID)
		protected.PUT("/maps/:id", mapHandler.UpdateMap)
		protected.DELETE("/maps/:id", mapHandler.DeleteMap)
	}

	return &SimpleServer{
		server:  router,
		authSrv: auth,
		mapSrv:  maps,
	}
}

func (s *SimpleServer) Run(ctx context.Context) {
	err := s.server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
