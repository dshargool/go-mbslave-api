package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/dshargool/go-mbslave-api.git/pkg/handlers"
	"github.com/dshargool/go-mbslave-api.git/pkg/types"
)

var config = types.Configuration{}

func main() {
	config_path := "config.json"
	slog.Info("Reading configuration file: " + config_path)

	config, err := config.ReadConfig(config_path)
	if err != nil {
		log.Fatal("Error reading config ", err)
	}
	myDb := types.SqlDb{}
	myDb.Open(config.DBPath)
	defer myDb.Close()

	myDb.CreateTable()
	myDb.UpdateTableTags(config.Registers)

	slog.Info("Starting modbus TCP slave")

	slog.Info("Starting handler")
	handler := handlers.New(config, &myDb)
	handler.MbSlave = handler.MbInit()
	handler.MbStart()
	defer handler.MbStop()
	handler.HandleRequests(config.Port)

	fmt.Println("End")
}
