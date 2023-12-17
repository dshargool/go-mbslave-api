package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/dshargool/go-mbslave-api.git/pkg/handlers"
	"github.com/dshargool/go-mbslave-api.git/pkg/types"
)

var config = types.Configuration{}

func main() {
	configPtr := flag.String("config", "config.json", "Application file path")
	flag.Parse()

	config_path := *configPtr
	slog.Info("Reading configuration file: " + config_path)

	config, err := config.ReadConfig(config_path)
	if err != nil {
		log.Fatal("Error reading config ", err)
	}
	myDb := types.SqlDb{}
	err = myDb.Open(config.DBPath)
	if err != nil {
		os.Exit(1)
	}
	defer myDb.Close()

	err = myDb.CreateTable()
	if err != nil {
		os.Exit(1)
	}
	myDb.UpdateTableTags(config.Registers)

	slog.Info("Starting modbus TCP slave")

	slog.Info("Starting handler")
	handler := handlers.New(config, &myDb)
	handler.MbSlave = handler.MbInit(config.ModbusPort)
	handler.MbStart()
	defer handler.MbStop()
	handler.HandleRequests(config.ApiPort)

	fmt.Println("End")
}
