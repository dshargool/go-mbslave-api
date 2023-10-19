# GO Modbus API
Simple API for translating a set of defined modbus registers from a configuration file into a web api.

## Configuration file
The configuration file has the format that outlines the equipment whose registers we want to translate.  The config file will translate into endpoints in the format:
`http://api-ip/<SLAVE_NAME>/<register_name>`

The fields available for configuration are:
| Field | Description |
| --- | --- |
| "name" | Name of the slave device |
| "slave_ip" | IP address of the Modbus TCP slave |
| "cache_ttl_s" | Cache time-to-live in seconds.  Used to reduce requests modbus endpoint and avoid overloading it |
| "registers:address" | The modbus register address |
| "registers:divisor" | The divisor for the data at the register address (allows decimals) |
| "registers:name"    | The name of the register that will be used to access the register data at the API |

```json
[
  {
    "name": "VFD_SLAVE_1",
    "slave_ip": "192.168.1.1",
    "cache_ttl_s": 60,
    "registers": [
      {
        "address": 40001,
        "divisor": 10,
        "name": "VFD2_Temp"
      },
      {
        "address": "40002",
        "divisor": 10,
        "name": "VFD2_Temp"
      }
    ],
    "string": "Hello World"
  },
  {
    "name": "VFD_SLAVE_2",
    "slave_ip": "192.168.1.2",
    "registers": [
      {
        "address": 40001,
        "divisor": 10,
        "name": "VFD2_Temp"
      },
      {
        "address": "40002",
        "divisor": 10,
        "name": "VFD2_Temp"
      }
    ],
    "string": "Hello World"
  }
]
```

With the data available in our configuration file we are able to make a variety of requests.

## API Requests

We can make requests to our endpoint using the configured endpoint and register names.  
(Optional idea: Allow requests by index (/0/0 would give first endpoints first register)

### GET
GET requests will retrieve the data from the API for the appropriately formed request.  This will work with all register types.

### PUT
PUT requests will only work for holding registers (to allow values to be written) and coils (to allow them to be set).  Other register types will not accept PUT requests.


