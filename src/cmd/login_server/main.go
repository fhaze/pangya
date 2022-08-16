package main

import (
	"pangya/src/internal/database"
	"pangya/src/internal/logger"
	"pangya/src/internal/sync"
	"pangya/src/internal/utils"
	loginserver "pangya/src/login_server"
	"pangya/src/login_server/handlers"
	syncclient "pangya/src/sync_client"
	synchandlers "pangya/src/sync_client/handlers"

	"github.com/common-nighthawk/go-figure"
	"github.com/joho/godotenv"
)

func main() {
	figure.NewColorFigure("LoginServer", "graffiti", "blue", true).Print()
	godotenv.Load()
	database.Connect()
	svc := loginserver.New()

	port := utils.GetIntEnv("LOGIN_PORT")
	client := syncclient.New(svc)

	svc.AddHandler(0x0001, handlers.NewP0x0001_ClientLogin(svc))

	client.AddHandler(sync.PacketHandshake, synchandlers.NewServerhandshake(client))
	client.AddHandler(sync.PacketGameServerList, synchandlers.NewServerGameServerList(client, svc))

	syncHost := utils.GetStringEnv("SYNC_HOST")
	syncPort := utils.GetIntEnv("SYNC_PORT")

	err := client.Dial(syncHost, syncPort)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	err = svc.Listen(port)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}
