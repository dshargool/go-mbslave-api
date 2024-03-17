package types

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

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
	if err := db.QueryRow("SELECT COUNT(name) FROM sqlite_master WHERE type='table' AND name='datapoints';").Scan(
		&exists); err != nil && err != sql.ErrNoRows {
		slog.Error("Failed to create new table", "error", err)
		return err
	}
	if !exists {
		results, err := db.Exec("CREATE TABLE datapoints (address VARCHAR(100) PRIMARY KEY NOT NULL, description VARCHAR(100), tag VARCHAR(75) NOT NULL, value REAL, datatype VARCHAR(10), last_update TEXT DEFAULT CURRENT_TIMESTAMP);")
		if err != nil {
			fmt.Println("failed to execute query", err)
			return err
		}
		slog.Info("Table created successfully", "result", results)
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
    RETURNING tag;`
	var err error
	for _, register := range registers {
		slog.Debug("Updating row", "reg", register)
		// Check to see if it's a multibit address.  If it is we create a generic one to r/w to
		if strings.Contains(register.Address, "_") {
            addr := strings.Split(register.Address, "_")[0]
            genReg := ModbusTag{
            	Tag:         "GenericAddressTag" + addr,
            	Description: "Generic Digital Address for " + addr,
            	Address:     addr,
            	DataType:    "digital",
            }
			err = db.QueryRow(queryStmt, &genReg.Address,
				&genReg.Description,
				&genReg.Tag, &genReg.DataType).Scan(&genReg.Tag)
            slog.Debug("Updating generic address table tag", "reg", genReg)
			if err != nil {
				slog.Error("failed to execute generic register query", "error", err)
				return
			}
		}
		err = db.QueryRow(queryStmt, &register.Address, &register.Description,
			&register.Tag, &register.DataType).Scan(&register.Tag)
        slog.Debug("Updating tag", "reg", register)
		if err != nil {
			slog.Error("failed to execute query", "error", err)
			return
		}
	}
}

func (db *SqlDb) GetRowByTag(tag string) (response ModbusResponse, err error) {
	slog.Debug("Getting DB Row", "tag", tag)
    addr, err := db.GetAddressByTag(tag)
    if err != nil {
        return response, err
    }
    return db.GetRowByAddress(addr)
}

func (db *SqlDb) GetAddressByTag(tag string) (response string, err error) {
    var resp ModbusResponse
	slog.Debug("Getting DB Row", "tag", tag)
	rows := db.QueryRow("SELECT address FROM datapoints WHERE tag=$1", tag)
	err = rows.Scan(&resp.Address)

	return resp.Address, err
}

func (db *SqlDb) SetTagValue(tag string, value float64) error {
	slog.Debug("Setting DB Row", "tag", tag, "value", value)
	_, err := db.Exec("UPDATE datapoints SET value = $1 WHERE tag = $2", value, tag)
        if err != nil {
            return err
        }
    dataType, err := db.GetDataTypeByTag(tag)
    if strings.Contains(dataType, "digital") {
        addr, err := db.GetAddressByTag(tag)
        if err != nil {
            return err
        }
        err = db.SetAddressValue(addr, value)
        if err != nil {
            return err
        }
    }
	if err != nil {
		return err
	}
	return nil
}

func (db *SqlDb) SetGenericBitAddress(address string, value float64) error{
        genAddress := strings.Split(address, "_")[0]
		digitShift, err := strconv.Atoi(strings.Split(address, "_")[1])
		if err != nil {
			return err
		}

        slog.Debug("Setting generic address", "addr", genAddress, "shift", digitShift, "value", value)
		currRow, err := db.GetRowByAddress(genAddress)
		currVal := uint64(currRow.Value)
		intVal := uint64(value)
		if err != nil && err == sql.ErrNoRows {
            slog.Error("FAILED TO GET ROW", "err", err, "row", currRow)
			return err
		} else if err != nil {
            currVal = 0
        }
        
        if intVal > 0 {
		    currVal |= 1 << uint64(digitShift)
        } else {
		    currVal &= 0 << uint64(digitShift)
        }

	    slog.Debug("Setting generic DB Row", "address", genAddress, "value", currVal)
		_, err = db.Exec("UPDATE datapoints SET value = $1 WHERE address = $2", currVal, genAddress)
		if err != nil {
			return err
		}
        return nil
}

func (db *SqlDb) GetGenericBitAddress(address string) (value int, err error) {
    splitStr := strings.Split(address, "_")
    genAddress := splitStr[0]
    if len(splitStr) == 1{
        return 0, errors.New("Length of address string does not contain digit information")
    }
	digitShift, err := strconv.Atoi(strings.Split(address, "_")[1])
	if err != nil {
		return 0, err
	}

    current, err := db.GetRowByAddress(genAddress)

    value = (int(current.Value) >> digitShift) & 1

    if err != nil {
        slog.Error("Could not find generic address", "addr", genAddress)
    }

    return value, nil
}

func (db *SqlDb) GetRowByAddress(address string) (response ModbusResponse, err error) {
	slog.Debug("Getting DB Row", "address", address)
	rows := db.QueryRow("SELECT address,tag,description,datatype,value,last_update FROM datapoints WHERE address=$1", address)
	err = rows.Scan(&response.Address, &response.Tag, &response.Description, &response.DataType, &response.Value, &response.LastUpdate)
    if err != nil && strings.Contains(err.Error(), "NULL to float64") && strings.Contains(response.DataType, "digital") && strings.Contains(response.Address, "_") {
        err = nil
        genValue, err := db.GetGenericBitAddress(response.Address)
        if err != nil {
            return response, err
        }
        response.Value = float64(genValue)
    } else {
        return response, err
    }
	return
}

func (db *SqlDb) GetDataTypeByTag(tag string) (dataType string, err error) {

	slog.Debug("Getting DB Row Datatype", "tag", tag)
	var db_dataType string = "none"
	rows := db.QueryRow("SELECT datatype FROM datapoints WHERE tag=$1", tag)
	err = rows.Scan(&db_dataType)

	return db_dataType, err
}

func (db *SqlDb) GetDataTypeByAddress(address string) (dataType string, err error) {
	slog.Debug("Getting DB Row Datatype", "address", address)
	var db_dataType string = "none"
	rows := db.QueryRow("SELECT datatype FROM datapoints WHERE address=$1", address)
	err = rows.Scan(&db_dataType)

	return db_dataType, err
}

func (db *SqlDb) SetAddressValue(address string, value float64) error {
	slog.Info("Setting DB Row", "address", address, "value", value)
	_, err := db.Exec("UPDATE datapoints SET value = $1 WHERE address = $2", value, address)
	if err != nil {
		return err
	}
    dataType, err := db.GetDataTypeByAddress(address)
    if err != nil {
        return err
    }
	// If we are sure this is a digital address
	if strings.Contains(dataType, "digital") && strings.Contains(address, "_") {
        _ = db.SetGenericBitAddress(address, value)
	}
	return nil
}

func (db *SqlDb) PropogateValueSubAddressDigital(address string) error {
    slog.Error("Propogation Nation for"+address)
    currentData, err := db.GetRowByAddress(address)
    if err != nil {
        return err
    }

    newValue := uint64(currentData.Value)

    rows, err := db.Query("SELECT address FROM datapoints WHERE address LIKE $1", address+"_%")
    if err != nil {
        slog.Error("No rows!")
        return err
    }
    defer rows.Close()

    var fullAddress string
    for rows.Next(){
        err := rows.Scan(&fullAddress)
        slog.Error("Got row", "addr", fullAddress)
        if err != nil {
            return err
        }

        digit, err := strconv.Atoi(strings.Split(fullAddress, "_")[1])
        if err != nil {
            return err
        }
        valToSet := newValue & 1 << digit
        err = db.SetAddressValue(fullAddress, float64(valToSet))
        if err != nil {
            return err
        }
    }
    slog.Error("Out of rows")
    return nil
}
