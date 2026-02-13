package main

import (
	"context"
	"log"

	"github.com/Oralkhan-coder/mind-map/config"
	internalHttp "github.com/Oralkhan-coder/mind-map/internal/http"
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

	server := internalHttp.NewSimpleServer()
	server.Run(ctx)
}
