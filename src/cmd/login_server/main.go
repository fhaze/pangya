package main

import (
	"pangya/src/internal/database"
	"pangya/src/internal/logger"
	"pangya/src/internal/sync"
	"pangya/src/internal/utils"
	loginserver "pangya/src/login_server"
	"pangya/src/login_server/handlers"
	syncclient "pangya/src/sync_client"

	"github.com/common-nighthawk/go-figure"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	figure.NewColorFigure("LoginServer", "graffiti", "blue", true).Print()
	godotenv.Load()
	database.Connect()
	svc := loginserver.New()

	svc.AddHandler(0x0001, handlers.NewP0x0001_ClientLogin())

	port := utils.GetIntEnv("LOGIN_PORT")
	client := syncclient.New()

	syncHost := utils.GetStringEnv("SYNC_HOST")
	syncPort := utils.GetIntEnv("SYNC_PORT")

	err := client.Dial(syncHost, syncPort)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	err = client.Handshake(sync.LoginServer)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	go synchronise(&client)

	err = svc.Listen(port)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func synchronise(client *syncclient.Client) {
	for {
		buff, err := (*client).Read()
		if err != nil {
			logger.Log.Fatal("error reading from sync server", zap.Error(err))
		}
		logger.Log.Debug("sync server read not implemented yet", zap.ByteString("buff", buff))
	}
}
