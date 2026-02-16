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
	nodeSrv NodeSrv
}

func NewSimpleServer(auth AuthSrv, maps MapSrv, node NodeSrv, cfg *config.SecretConfig) *SimpleServer {
	authHandler := transport.NewAuthHandler(auth)
	mapHandler := transport.NewMapHandler(maps)
	nodeHandler := transport.NewNodeHandler(node)

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

		protected.GET("/maps/:id/nodes", nodeHandler.GetByMapId)
		protected.POST("/maps/:id/nodes", nodeHandler.CreateNode)
	}

	return &SimpleServer{
		server:  router,
		authSrv: auth,
		mapSrv:  maps,
		nodeSrv: node,
	}
}

func (s *SimpleServer) Run(ctx context.Context) {
	err := s.server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
