**Lenpaste** is a web service that allows you to share notes anonymously, an alternative to `pastebin.com`.


## Features
- No need to register
- Uses cookies only to store settings
- Can work without JavaScript
- Has its own API
- Open source and self-hosted



## Public servers list
| Server                                             | Description                                                                                                       |
| ---------------------------------------------------| ----------------------------------------------------------------------------------------------------------------- |
| [paste.lcomrade.su](https://paste.lcomrade.su)     | Server managed by the Lenpaste developer.                                                                         |
| [code.dbt3ch.com](https://code.dbt3ch.com)         | Server is managed by DB Tech. He made a [video about Lenpaste v1.1](https://www.youtube.com/watch?v=YxcHxsZHh9A). |
| [notepad.co.il](https://notepad.co.il)             | Server managed by Shlomi Porush. He reported the bug and made some suggestions.                                   |
| [lenp.pardesicat.xyz](https://lenp.pardesicat.xyz) | Server managed by Pardesi_Cat. He helped correct the documentation.                                               |

Find more public servers here or add your own: https://monitor.lcomrade.su/?srv=lenpaste



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
	# There are images for x86, x64, ARM64, ARM v7, ARM v6.
	# The Raspberry Pi is supported, including the latest 64-bit versions.
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
      - LENPASTE_BODY_MAX_LENGTH=20000        # Maximum length of the paste body. If -1 disable length limit. Can't be -1.
      - LENPASTE_MAX_PASTE_LIFETIME=unlimited # Maximum lifetime of the paste. Examples: 10m, 1h 30m, 12h, 7w, 30d, 365d.
      - LENPASTE_NEW_PASTES_PER_5MIN=15       # Maximum number of paste that can be created in 5 minutes from one IP. If 0 disable rate-limit.
      - LENPASTE_ADMIN_NAME=                  # Name of the administrator of this server.
      - LENPASTE_ADMIN_MAIL=                  # Email of the administrator of this server.
      - LENPASTE_UI_DEFAULT_LIFETIME=         # Lifetime of paste will be set by default in WEB interface. Examples: 10min, 1h, 1d, 2w, 6mon, 1y.
    volumes:
      # /data/lenpaste.db - SQLite DB if used.
      # /data/about       - About this server (TXT file).
      # /data/rules       - This server rules (TXT file).
      # /data/terms       - This server "terms of use" (TXT file).
      # /data/lenpasswd   - If this file exists, the server will ask for auth to create new pastes.
      #                     File format: USER:PLAIN_PASSWORD on each line.
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
An official image may not support your architecture e.g. MIPS, PowerPC, etc.
So you can build your own image to run on an officially unsupported architecture
(of course you have to rebuild it every time you update Lenpaste).

On Debian/Ubuntu:
```
sudo apt update
sudo apt -y install git docker docker.io
git clone https://git.lcomrade.su/root/lenpaste.git
cd ./lenpaste/
git checkout vX.X
sudo docker build -t localhost/lenpaste:latest ./
```

The `localhost/lenpaste:latest` image should now have appeared on your local machine.
You can use it in `docker-compose.yml` or copy it to another machine.



## Other documentation
For all:
- [Roadmap for new release](ROADMAP.md)
- [Frequently Asked Questions](FAQ.md)

For instance administrators:
- [Reverse proxy: Nginx](docs/reverse_proxy_nginx.md)
- [Database: PostgreSQL](docs/db_postgresql.md)
- [Rate limiting](docs/ratelimits.md)
- [Make Lenpaste server private](docs/private_server.md)

For developers:
- [Lenpaste API](https://paste.lcomrade.su/docs/apiv1)
- [Libraries for working with API](https://paste.lcomrade.su/docs/api_libs)



## Might be interesting
Manuals:
- [How to Install LenPaste on Your Synology NAS](https://mariushosting.com/how-to-install-lenpaste-on-your-synology-nas/) (WEB site)
- [Lenpaste | TrueCharts](https://truecharts.org/docs/charts/incubator/lenpaste/) (WEB site)

Reviews:
- [Pastebin Clone in Docker with Lenpaste](https://www.youtube.com/watch?v=YxcHxsZHh9A) (YouTube video)



## Bugs and Suggestion
If you have any questions or suggestions, you can write here:
- Join to Matrix room: [`#lenpaste:lcomrade.su`](https://matrix.to/#/#lenpaste:lcomrade.su)
- Contact me: Leonid Maslakov <root@lcomrade.su>



## Donate
All donations will go to Leonid Maslakov, for now the sole developer:
- Qiwi: https://qiwi.com/n/LCOMRADE
