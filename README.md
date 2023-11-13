# GO Modbus Slave API
API to allow reads and writes to a central modbus slave to occur based on the OPC tag we want to interact with.  This allows both the API and an external modbus master to query the modbus slave to read/write data.

## Configuration file
The configuration file has the format that outlines the equipment whose registers we want to translate.  The config file will translate into endpoints in the format:
`http://api-ip/tag/registers:tag`

The fields available for configuration are:
| Field | Description |
| --- | --- |
| "port" | Port for API access |
| "db" | Path to sqlite database |
| "registers:tag" | OPC/Instrument Tag to access this data point via API |
| "registers:name"    | The name of the register that will be used to access the register data at the API |
| "registers:address" | The modbus register address |
| "registers:divisor" | The divisor for the data at the register address (allows decimals) |

```json
{
    "port": 8081,
    "db": "./db/test.db",
    "registers": [
        {
            "tag": "210XT1055.PNT",
            "name": "TestPoint1",
            "address": 40001,
            "divisor": 10,
        },
        {
            "tag": "210XT2055.PNT",
            "name": "TestPoint2",
            "address": 40002,
            "divisor": 20,
        }
    ]
}
```

With the data available in our configuration file we are able to make a variety of requests.

## API Requests

We can make requests to our endpoint using the configured endpoint and register names.  
(Optional idea: Allow requests by index (/InputRegister/0 would give first endpoints first input register)

### GET
GET requests will retrieve the data for the requested appropriate data point

### PUT
PUT requests allow data to be written to any of the data points.

## MODBUS Requests

We can make modbus requests to our endpoint using the configured endpoint and register addresses.  This application acts as the modbus slave so only responds to requests and will not make them on its own.

## Database
The database stores the current data points; this allows us to consistently reboot the application without losing the state that needs to be transfered.  This means that our database values should be as close to the most recent ones from either the API or Modbus Master to be communicated.

### Tables
There is a single main table for our data points.  The register address acts as our primary key.
TABLE: tags
Columns: address, divisor, value 
