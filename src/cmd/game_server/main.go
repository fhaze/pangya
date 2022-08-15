package main

import (
	gameserver "pangya/src/game_server"
	"pangya/src/internal/database"
	"pangya/src/internal/logger"
	"pangya/src/internal/sync"
	"pangya/src/internal/utils"
	syncclient "pangya/src/sync_client"
	synchandlers "pangya/src/sync_client/handlers"

	"github.com/common-nighthawk/go-figure"
	"github.com/joho/godotenv"
)

func main() {
	figure.NewColorFigure("GameServer", "graffiti", "green", true).Print()
	godotenv.Load()
	database.Connect()
	svc := gameserver.New()

	port := utils.GetIntEnv("GAME_PORT")
	client := syncclient.New(svc)

	client.AddHandler(sync.PacketHandshake, synchandlers.NewServerhandshake(client))

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
