package main

import (
	"pangya/src/internal/database"
	loginserver "pangya/src/login_server"
	"pangya/src/login_server/handlers"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()
	database.Connect()
	svc := loginserver.New()
	svc.AddHandler(0x0001, handlers.NewP0x0001_ClientLogin())
	err := svc.Listen(10103)
	if err != nil {
		zap.Error(err)
	}
}
