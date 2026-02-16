package main

import (
	"context"
	"log"

	"github.com/Oralkhan-coder/mind-map/config"
	internalHttp "github.com/Oralkhan-coder/mind-map/internal/http"
	"github.com/Oralkhan-coder/mind-map/internal/service"
	mongo "github.com/Oralkhan-coder/mind-map/pkg"
)

func main() {
	ctx := context.Background()
	cfg := config.InitConfig()

	db, err := mongo.NewDB(ctx, cfg.Mongo)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Close(ctx)

	emailService := service.NewEmailService(cfg.Email)
	userService := service.NewUserService(db.Database.Collection("users"))
	authService := service.NewAuthService(userService, emailService, cfg.Secret)
	mapService := service.NewMapService(db.Database.Collection("maps"), userService)

	server := internalHttp.NewSimpleServer(authService, mapService, cfg.Secret)
	server.Run(ctx)
}
