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
	authService := service.NewAuthService(db.Database.Collection("users"), emailService, cfg.Secret)

	server := internalHttp.NewSimpleServer(authService, cfg.Secret)
	server.Run(ctx)
}
