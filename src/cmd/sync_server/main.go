package main

import (
	"os"
	"pangya/src/internal/logger"
	"pangya/src/internal/sync"
	syncserver "pangya/src/sync_server"
	"pangya/src/sync_server/handlers"
	"strconv"

	"github.com/common-nighthawk/go-figure"
	"github.com/joho/godotenv"
)

func main() {
	figure.NewColorFigure("SyncServer", "graffiti", "red", true).Print()
	godotenv.Load()
	svc := syncserver.New()

	svc.AddHandler(sync.PacketHandshake, handlers.NewClientHandshake(svc))

	port, err := strconv.Atoi(os.Getenv("SYNC_PORT"))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	if port == 0 {
		logger.Log.Fatal("port cannot be 0")
	}

	err = svc.Listen(port)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}
