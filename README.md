**Lenpaste** is a web service that allows you to share notes anonymously, an alternative to `pastebin.com`.


## Features
- No need to register
- Uses cookies only to store settings
- Can work without JavaScript
- Has its own API
- Open source and self-hosted

Find out what's coming in the next release on the [roadmap](ROADMAP.md).



## Public servers list
| Server                                         | Description                                                                                                       |
| ---------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| [paste.lcomrade.su](https://paste.lcomrade.su) | Server managed by the Lenpaste developer.                                                                         |
| [code.dbt3ch.com](https://code.dbt3ch.com)     | Server is managed by DB Tech. He made a [video about Lenpaste v1.1](https://www.youtube.com/watch?v=YxcHxsZHh9A). |



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
  	# If you want to run Lenpaste on your Raspberry Pi, use 'latest-armhf' instead of 'latest'.
    image: git.lcomrade.su/root/lenpaste:latest
    restart: always
    environment:
      # All parameters are optional
      - LENPASTE_ADDRESS=:80                  # ADDRES:PORT for HTTP server.
      - LENPASTE_DB_DRIVER=sqlite3            # Currently supported drivers: 'sqlite3' and 'postgres'.
      - LENPASTE_DB_SOURCE=/data/lenpaste.db  # DB source.
      - LENPASTE_DB_CLEANUP_PERIOD=3h         # Interval at which the DB is cleared of expired but not yet deleted pastes.
      - LENPASTE_ROBOTS_DISALLOW=false        # Prohibits search engine crawlers from indexing site using robots.txt file.
      - LENPASTE_TITLE_MAX_LENGTH=100         # Maximum length of the paste title. If 0 disable title, if -1 disable length limit.
      - LENPASTE_BODY_MAX_LENGTH=10000        # Maximum length of the paste body. If -1 disable length limit. Can't be -1.
      - LENPASTE_MAX_PASTE_LIFETIME=unlimited # Maximum lifetime of the paste. Examples: 10m, 1h 30m, 12h, 7w, 30d, 365d.
      - LENPASTE_ADMIN_NAME=                  # Name of the administrator of this server.
      - LENPASTE_ADMIN_MAIL=                  # Email of the administrator of this server.
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

TIP: If you want to install updates, run: `docker-compose pull && docker-compose up -d && docker system prune -a -f`



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



## Build Docker image
**Why is it necessary?**
An official image may not support your architecture e.g. ARM64, PowerPC, etc.
So you can build your own image to run on an officially unsupported architecture
(of course you have to rebuild it every time you update Lenpaste).

On Debian/Ubuntu:
```
sudo apt update
sudo apt -y install git docker docker.io
git clone https://git.lcomrade.su/root/lenpaste.git
cd ./lenpaste/
git checkout vX.X
sudo docker -t localhost/lenpaste:latest ./
```

The `localhost/lenpaste:latest` image should now have appeared on your local machine.
You can use it in `docker-compose.yml` or copy it to another machine.



## Other documentation
For all:
- [Roadmap for new release](ROADMAP.md)

For instance administrators:
- [Reverse proxy: Nginx](docs/reverse_proxy_nginx.md)
- [Database: PostgreSQL](docs/db_postgresql.md)

For developers:
- [Lenpaste API](https://paste.lcomrade.su/docs/apiv1)
- [Libraries for working with API](https://paste.lcomrade.su/docs/api_libs)



## Might be interesting
Reviews and testimonials:
- [Pastebin Clone in Docker with Lenpaste](https://www.youtube.com/watch?v=YxcHxsZHh9A) (YouTube video)



## Bugs and Suggestion
If you want to:
- Report a bug.
- Ask a question.
- Become a contributor.
- Or something else.

Then write to me: Leonid Maslakov <root@lcomrade.su>
