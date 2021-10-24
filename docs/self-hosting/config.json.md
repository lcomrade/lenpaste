# config.json
This is the main configuration file (format: JSON).
Here is an example of a configuration file:

```
{
	"HTTP": {
		"Listen": ":8000",
		"UseTLS": false,
		"SSLCert": "./data/fullchain.pem",
		"SSLKey": "./data/privkey.pem"
	},
	"Storage": {
		"CleanJobPeriod": 600
	},
	"Logs": {
		"SaveErr": true,
		"SaveInfo": true,
		"SaveJob": true,
		"RotateLogs": true,
		"MaxLogSize": 1000000
	}
}
```

## HTTP
### Listen
Specifies the port and IP address from which the web server will receive requests. Examples:

`:80` - Accept requests from all IP addresses on port 80

`127.0.0.1:8000` - Receive requests from localhost on port 8000

### UseTLS
If `true` enables TLS encryption. Requires `SSLCert` and `SSLKey`.

### SSLCert
Specifies the path to the SSL certificate.

### SSLKey
Specifies the path to the SSL private key.


## Storage
### CleanJobPeriod
Specifies the frequency of background deletion of expired pastes.
Specifies in seconds.


## Logs
### SaveErr
Responsible for saving logs marked `ERROR` to file `./data/log/error`.

### SaveInfo
Responsible for saving logs marked `INFO` to file `./data/log/info`.

### SaveJob
Responsible for saving logs marked `JOB` to file `./data/log/job`.

### RotateLogs
Responsible for rotating the logs.
Rotation is performed when the app is started.

### MaxLogSize
Maximum log size in bytes. Used when rotating logs.
