# GO Modbus Slave API
API to allow reads and writes to a central modbus slave to occur based on the OPC tag we want to interact with.  This allows both the API and an external modbus master to query the modbus slave to read/write data.

## Configuration file
The configuration file has the format that outlines the equipment whose registers we want to translate.  The config file will translate into endpoints in the format:
`http://api-ip/tag/registers:tag`

The fields available for configuration are:
| Field | Description |
| --- | --- |
| "api_port" | Port for API access |
| "modbus_port" | Port for Modbus Slave access |
| "db" | Path to sqlite database |
| "allow_null_registers" | Allow reading of registers that aren't configured |
| "registers:tag" | API Tag to access this data point via API |
| "registers:name"    | The name of the register that will be used to access the register data at the API |
| "registers:address" | The modbus holding register address |
| "registers:datatype" | The datatype stored at the register address (will read multiple if datatype size is larger than 16 bits |

```json
{
    "api_port": 8081,
    "modbus_port": 6502,
    "db": "./db/test.db",
    "allow_null_registers": true,
    "registers": [
        {
            "tag": "TestTag1",
            "name": "TestPoint1",
            "address": 40001,
            "datatype": "float32"
        },
        {
            "tag": "TestTag2",
            "name": "TestPoint2",
            "address": 40003,
            "datatype": "float32"
        }
    ]
}
```

With the data available in our configuration file we are able to make a variety of requests.

## API Requests

We can make requests to our endpoint using the configured endpoint and register names.

Valid endpoints are `*/tag/<tag>` and `*/register/<address>` where both `<tag>` and `<address>` are from the configuration file.

### GET

GET requests will retrieve the data for the requested appropriate data point

### PUT

PUT requests allow data to be written to any of the data points.

## MODBUS Requests

We can make modbus requests to our endpoint using the configured endpoint and register addresses.  This application acts as the modbus slave so only responds to requests and will not make them on its own.

## Data Types

Support for basic datatypes are available; `float32`, `float64`, `int16`, `uint16`.  Unsupported datatypes will return an error.

Floats are assembled in the byte sequence of upper byte then lower byte with big endianness.

## Database

The database stores the current data points; this allows us to consistently reboot the application without losing the state that needs to be transfered.  This means that our database values should be as close to the most recent ones from either the API or Modbus Master to be communicated.

### Tables

There is a single main table for our data points.  The register address acts as our primary key.
TABLE: datapoints
Columns: address, description, datatype, value, last_updated
