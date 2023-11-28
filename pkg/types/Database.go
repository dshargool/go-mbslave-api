package types

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

func (db *SqlDb) Open(dbPath string) {
	slog.Info("Opening sqlite3 database at: " + dbPath)
	newDb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	db.DB = newDb
}

func (db *SqlDb) CreateTable() {
	var exists bool
	if err := db.QueryRow("SELECT COUNT(name) FROM sqlite_master WHERE type='table' AND name='datapoints';").Scan(&exists); err != nil && err != sql.ErrNoRows {
		fmt.Println("failed to execute exists query", err)
		return
	}
	if !exists {
		results, err := db.Exec("CREATE TABLE datapoints (address INTEGER PRIMARY KEY NOT NULL, description VARCHAR(100), tag VARCHAR(75) NOT NULL, value REAL, divisor INTEGER, last_update TEXT DEFAULT CURRENT_TIMESTAMP);")
		if err != nil {
			fmt.Println("failed to execute query", err)
			return
		}
		slog.Info("Table created successfully", results)
	} else {
		slog.Info("Table 'datapoints' already exists ")
	}
}

func (db *SqlDb) UpdateTableTags(registers map[OpcTag]ModbusTag) {
	queryStmt := `INSERT INTO datapoints (address,description,tag,divisor) VALUES
    ($1, $2, $3, $4) 
    ON CONFLICT(address) DO UPDATE SET
    description=excluded.description, tag=excluded.tag, divisor=excluded.divisor
    RETURNING address;`
	for _, register := range registers {
		err := db.QueryRow(queryStmt, &register.Address, &register.Description, &register.Tag, &register.Divisor).Scan(&register.Address)
		if err != nil {
			slog.Error("failed to execute query", "error", err)
			return
		}
	}
}

func (db *SqlDb) GetRowByTag(tag string) (response ModbusResponse, err error) {
	slog.Info("Getting DB Row", "tag", tag)
	rows := db.QueryRow("SELECT address,tag,description,divisor,value,last_update FROM datapoints WHERE tag=$1", tag)
	err = rows.Scan(&response.Address, &response.Tag, &response.Description, &response.Divisor, &response.Value, &response.LastUpdate)

	return
}

func (db *SqlDb) SetTagValue(tag string, value float64) error {
	slog.Info("Setting DB Row", "tag", tag, "value", value)
	_, err := db.Exec("UPDATE datapoints SET value = $1 WHERE tag = $2", value, tag)
	if err != nil {
		return err
	}
	return nil
}

func (db *SqlDb) GetRowByAddress(address int) (response ModbusResponse, err error) {
	slog.Info("Getting DB Row", "address", address)
	rows := db.QueryRow("SELECT address,tag,description,divisor,value,last_update FROM datapoints WHERE address=$1", address)
	err = rows.Scan(&response.Address, &response.Tag, &response.Description, &response.Divisor, &response.Value, &response.LastUpdate)

	return
}

func (db *SqlDb) SetAddressValue(address int, value float64) error {
	slog.Info("Setting DB Row", "address", address, "value", value)
	_, err := db.Exec("UPDATE datapoints SET value = $1 WHERE address = $2", value, address)
	if err != nil {
		return err
	}
	return nil
}
