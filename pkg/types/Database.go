package types

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

func (db *SqlDb) Open(dbPath string) error {
	slog.Info("Opening sqlite3 database at: " + dbPath)
	newDb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.Error("Could not open sqlite3 db", "error", err.Error())
		return err
	}
	db.DB = newDb
	return nil
}

func (db *SqlDb) CreateTable() error {
	var exists bool
	if err := db.QueryRow("SELECT COUNT(name) FROM sqlite_master WHERE type='table' AND name='datapoints';").Scan(&exists); err != nil && err != sql.ErrNoRows {
		slog.Error("Failed to create new table", "error", err)
		return err
	}
	if !exists {
		results, err := db.Exec("CREATE TABLE datapoints (address INTEGER PRIMARY KEY NOT NULL, description VARCHAR(100), tag VARCHAR(75) NOT NULL, value REAL, datatype VARCHAR(10), last_update TEXT DEFAULT CURRENT_TIMESTAMP);")
		if err != nil {
			fmt.Println("failed to execute query", err)
			return err
		}
		slog.Info("Table created successfully", results)
	} else {
		slog.Info("Table 'datapoints' already exists ")
	}
	return nil
}

func (db *SqlDb) UpdateTableTags(registers map[InstrumentTag]ModbusTag) {
	queryStmt := `INSERT INTO datapoints (address,description,tag,datatype) VALUES
    ($1, $2, $3, $4) 
    ON CONFLICT(address) DO UPDATE SET
    description=excluded.description, tag=excluded.tag, datatype=excluded.datatype
    RETURNING address;`
	for _, register := range registers {
		err := db.QueryRow(queryStmt, &register.Address, &register.Description, &register.Tag, &register.DataType).Scan(&register.Address)
		if err != nil {
			slog.Error("failed to execute query", "error", err)
			return
		}
	}
}

func (db *SqlDb) GetRowByTag(tag string) (response ModbusResponse, err error) {
	slog.Info("Getting DB Row", "tag", tag)
	rows := db.QueryRow("SELECT address,tag,description,datatype,value,last_update FROM datapoints WHERE tag=$1", tag)
	err = rows.Scan(&response.Address, &response.Tag, &response.Description, &response.DataType, &response.Value, &response.LastUpdate)

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
	rows := db.QueryRow("SELECT address,tag,description,datatype,value,last_update FROM datapoints WHERE address=$1", address)
	err = rows.Scan(&response.Address, &response.Tag, &response.Description, &response.DataType, &response.Value, &response.LastUpdate)

	return
}

func (db *SqlDb) GetDataTypeByAddress(address int) (dataType string, err error) {
	slog.Info("Getting DB Row Datatype", "address", address)
	var db_dataType string
	rows := db.QueryRow("SELECT datatype FROM datapoints WHERE address=$1", address)
	err = rows.Scan(&db_dataType)

	return db_dataType, err
}

func (db *SqlDb) SetAddressValue(address int, value float64) error {
	slog.Info("Setting DB Row", "address", address, "value", value)
	_, err := db.Exec("UPDATE datapoints SET value = $1 WHERE address = $2", value, address)
	if err != nil {
		return err
	}
	return nil
}
