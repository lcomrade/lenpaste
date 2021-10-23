# Configuration
Location of configuration files:
- Outside of Docker container: `/var/lib/lenpaste/`
- Inside Docker container: `/app/data/`

## config.json
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
	"Logs": {
		"SaveErr": true,
		"SaveInfo": true,
		"SaveJob": true,
		"RotateLogs": true,
		"MaxLogSize": 1000000
	}
}
```

## about.txt
Here you can write information about your server.
It will be displayed on the main page.

## rules.txt
In this file you can describe the rules for using your server.
It will be available at `/rules`.

## robots.txt
As the name suggests, this file is a `robots.txt` file.
If this file does not exist, the following contents will be used:

```
User-agent: *
Disallow: /
```
