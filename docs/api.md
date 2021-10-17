# Lenpaste API
Lenpaste has an open API. It does not require registration.

## Methods
### /api/new
Create a new paste.

Request:
 - `text` - Paste text
 - `expiration` - Paste expiration date
 
   | Value | Description |
   | ----- | ----------- |
   | 5m    | 5 minutes   |
   | 10m   | 10 minutes  |
   | 20m   | 20 minutes  |
   | 30m   | 30 minutes  |
   | 40m   | 40 minutes  |
   | 50m   | 50 minutes  |
   | 1h    | 1 hour      |
   | 2h    | 2 hour      |
   | 4h    | 4 hour      |
   | 12h   | 12 hour     |
   | 1d    | 1 day       |
   | 2d    | 2 day       |
   | 3d    | 3 day       |
   | 4d    | 4 day       |
   | 5d    | 5 day       |
   | 6d    | 6 day       |
   | 1w    | 1 week      |
   | 2w    | 2 week      |
   | 3w    | 3 week      |

Response: JSON

Response example:
```
{
	"Name":"8HslINRp"
}
```

### /api/get/PASTE_NAME
Get information about the paste and its text.

Request: NONE

Response: JSON

Response example:
```
{
	"Name":"8HslINRp",
	"Text":"My paste\r\nEND",
	"Info":{
		"CreateTime":1633845370,
		"DeleteTime":1633847170
	}
}
```

### /api/about
Get information about this server.

Request: NONE

Response: JSON

Response example:
```
{
	"Exist":true,
	"Text":"ABOUT\n1\n2\n3\n"
}
```

### /api/rules
Get information about server rules.

Request: NONE

Response: JSON

Response example:
```
{
	"Exist":true,
	"Text":"My server\nMy rules"
}
```

### /api/version
Get information about server version.

Request: NONE

Response: JSON

Response example:
```
{
	"Version":"v0.1-stable",
	"GitTag":"v0.1",
	"GitCommit":"30bb0a7e768f46ab02b74cf05fb948e8696a055c",
	"BuildDate":"16.10.2021"
}
```
