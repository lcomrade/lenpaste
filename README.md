**Lenpaste** is a web service that allows you to share notes anonymously, an alternative to `pastebin.com`.


## Features
- No need to register
- Does not use cookies
- Can work without JavaScript
- Has its own API
- Open source and self-hosted



## Public servers list
| Server                                         | Description                              |
| ---------------------------------------------- | ---------------------------------------- |
| [paste.lcomrade.su](https://paste.lcomrade.su) | Server managed by the Lenpaste developer |



## Launch your own server
1. If you don't already have Docker installed, do so:
```
apt-get install -y docker docker.io docker-compose
```

2. Use a file like this `docker-compose.yml`:
```yaml
version: "2"

services:
  lenpaste:
    image: git.lcomrade.su/root/lenpaste:latest
    restart: always
    environment:
      # All parameters are optional
      - LENPASTE_ADDRESS=:80                 # Set -address flag
      - LENPASTE_DB_DRIVER=sqlite3           # Set -db-driver flag
      - LENPASTE_DB_SOURCE=/data/lenpaste.db # Set -db-source flag
      - LENPASTE_DB_CLEANUP_PERIOD=3h        # Set -db-cleanup-period flag
      - LENPASTE_ROBOTS_DISALLOW=false       # If true set -robots-disallow flag
      - LENPASTE_TITLE_MAX_LENGTH=100        # Set -title-max-length flag. If 0 disable title, if -1 disable length limit.
      - LENPASTE_BODY_MAX_LENGTH=10000       # Set -body-max-length flag. If -1 disable length limit. Can't be -1.
      - LENPASTE_MAX_PASTE_LIFETIME=never    # Set -max-paste-lifetime flag. Examples: 2d, 12h, 7m.
      - LENPASTE_ADMIN_NAME=                 # Set -admin-name flag.
      - LENPASTE_ADMIN_MAIL=                 # Set -admin-mail flag.
    volumes:
      # /data/lenpaste.db - SQLite DB
      # /data/about.html  - About this server
      # /data/rules.html  - This server rules
      - "${PWD}/data:/data"
      - "/etc/timezone:/etc/timezone:ro"
      - "/etc/localtime:/etc/localtime:ro"
    ports:
      - "80:80"
```

3. Execute while in the directory where `docker-compose.yml` is located:
```
docker-compose pull && docker-compose up -d
```

PS: If you want to install updates, run: `docker-compose pull && docker-compose stop && docker-compose up -d && docker system prune -a -f`



## Build from source code
On Debian/Ubuntu:
```
sudo apt update
sudo apt -y install git make gcc golang
git clone https://git.lcomrade.su/root/lenpaste.git
cd ./lenpaste/
make
```

You can find the result of the build in the `./dist/` directory.



## Other documentation
For instance administrators:
- [Reverse proxy: Nginx](docs/reverse_proxy_nginx.md)
- [Database: PostgreSQL](docs/db_postgresql.md)

For developers:
- [Lenpaste API](https://paste.lcomrade.su/docs/apiv1)
- [Libraries for working with API](https://paste.lcomrade.su/docs/api_libs)



## Bugs and Suggestion
If you find a bug or have a suggestion, email me at root@lcomrade.su.
