package main

import (
	loginserver "pangya/login_server"
	"pangya/login_server/handlers"

	"go.uber.org/zap"
)

func main() {
	svc := loginserver.New()
	svc.AddHandler(0x0001, handlers.NewP0x0001_ClientLogin())
	err := svc.Listen(10103)
	if err != nil {
		zap.Error(err)
	}
}
