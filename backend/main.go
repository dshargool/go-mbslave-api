package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/dshargool/go-mbslave-api.git/pkg/handlers"
	"github.com/dshargool/go-mbslave-api.git/pkg/types"
)

var config = types.Configuration{}

func main() {
	configPtr := flag.String("config", "config.json", "Config File Path")
	dbPtr := flag.String("database", "", "Database File Path")
	flag.Parse()

	config_path := *configPtr
	slog.Info("Reading configuration file: " + config_path)

	config, err := config.ReadConfig(config_path)
	if err != nil {
		log.Fatal("Error reading config ", err)
	}

	// Set log level based on configuration
	switch strings.ToLower(config.LogLevel) {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "warn", "warning":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo) // Default to info level
	}

	db_path := *dbPtr

	if db_path == "" {
		db_path = config.DBPath
	}
	slog.Info("Opening database file: " + db_path)

	myDb := types.SqlDb{}
	err = myDb.Open(db_path)
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
