package types

import (
	"encoding/json"
	"log/slog"
	"os"
)

type ConfigurationData struct {
	ApiPort           int         `json:"api_port"`
	ModbusPort        int         `json:"modbus_port"`
	DBPath            string      `json:"db"`
	Description       string      `json:"description"`
	Registers         []ModbusTag `json:"registers"`
	AllowNullRegister bool        `json:"allow_null_register"`
}
type Configuration struct {
	ApiPort           int
	ModbusPort        int
	DBPath            string
	AllowNullRegister bool
	Registers         map[InstrumentTag]ModbusTag
}

func (c Configuration) ReadConfig(fileName string) (Configuration, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return Configuration{}, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configData := ConfigurationData{}
	err = decoder.Decode(&configData)
	if err != nil {
		return Configuration{}, err
	}
	config := configData.dataToConfiguration()
	slog.Info("Configuration found", "config", config)
	return config, nil
}

func (c ConfigurationData) dataToConfiguration() Configuration {
	config := Configuration{}
	config.Registers = make(map[InstrumentTag]ModbusTag)
	config.ApiPort = c.ApiPort
	config.ModbusPort = c.ModbusPort
	config.DBPath = c.DBPath
	config.AllowNullRegister = c.AllowNullRegister
	for _, reg := range c.Registers {
		config.Registers[InstrumentTag(reg.Tag)] = reg
	}
	return config
}
