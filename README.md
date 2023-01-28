**Lenpaste** is a web service that allows you to share notes anonymously, an alternative to `pastebin.com`.


## Features
- No need to register
- Supports multiple languages
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
| [lenp.pardesicat.xyz](https://lenp.pardesicat.xyz) | Server managed by Pardesi_Cat. He translated Lenpaste into Bengali and helped correct the documentation.          |

Find more public servers here or add your own: https://monitor.lcomrade.su/?srv=lenpaste



## Launch your own server
1. If you don't already have Docker installed, do so:
```
apt-get install -y docker.io docker-compose
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
      #
      # HTTP server
      - LENPASTE_ADDRESS=:80                  # ADDRES:PORT for HTTP server.
      #
      # Database settings
      - LENPASTE_DB_DRIVER=sqlite3            # Currently supported drivers: 'sqlite3' and 'postgres'.
      - LENPASTE_DB_SOURCE=/data/lenpaste.db  # DB source.
      - LENPASTE_DB_MAX_OPEN_CONNS=25         # Maximum number of connections to the database.
      - LENPASTE_DB_MAX_IDLE_CONNS=5          # Maximum number of idle connections to the database.
      - LENPASTE_DB_CLEANUP_PERIOD=3h         # Interval at which the DB is cleared of expired but not yet deleted pastes.
      #
      # Search engines
      - LENPASTE_ROBOTS_DISALLOW=false        # Prohibits search engine crawlers from indexing site using robots.txt file.
      #
      # Storage limits
      - LENPASTE_TITLE_MAX_LENGTH=100         # Maximum length of the paste title. If 0 disable title, if -1 disable length limit.
      - LENPASTE_BODY_MAX_LENGTH=20000        # Maximum length of the paste body. If -1 disable length limit. Can't be -1.
      - LENPASTE_MAX_PASTE_LIFETIME=unlimited # Maximum lifetime of the paste. Examples: 10m, 1h 30m, 12h, 7w, 30d, 365d.
      #
      # Rate limits
      - LENPASTE_GET_PASTES_PER_5MIN=50       # Maximum number of pastes that can be VIEWED in 5 minutes from one IP. If 0 disable rate-limit.
      - LENPASTE_GET_PASTES_PER_15MIN=100     # Maximum number of pastes that can be VIEWED in 15 minutes from one IP. If 0 disable rate-limit.
      - LENPASTE_GET_PASTES_PER_1HOUR=500     # Maximum number of pastes that can be VIEWED in 1 hour from one IP. If 0 disable rate-limit.
      - LENPASTE_NEW_PASTES_PER_5MIN=15       # Maximum number of pastes that can be CREATED in 5 minutes from one IP. If 0 disable rate-limit.
      - LENPASTE_NEW_PASTES_PER_15MIN=30      # Maximum number of pastes that can be CREATED in 15 minutes from one IP. If 0 disable rate-limit.
      - LENPASTE_NEW_PASTES_PER_1HOUR=40      # Maximum number of pastes that can be CREATED in 1 hour from one IP. If 0 disable rate-limit.
      #
      # Information about server admin
      - LENPASTE_ADMIN_NAME=                  # Name of the administrator of this server.
      - LENPASTE_ADMIN_MAIL=                  # Email of the administrator of this server.
      #
      # WEB interface settings
      - LENPASTE_UI_DEFAULT_LIFETIME=         # Lifetime of paste will be set by default in WEB interface. Examples: 10min, 1h, 1d, 2w, 6mon, 1y.
      - LENPASTE_UI_DEFAULT_THEME=dark        # Sets the default theme for the WEB interface. Examples: dark, light.
    volumes:
      # /data/lenpaste.db - SQLite DB if used.
      # /data/about       - About this server (TXT file).
      # /data/rules       - This server rules (TXT file).
      # /data/terms       - This server "terms of use" (TXT file).
      # /data/themes/*    - External WEB interface themes.
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
git checkout vX.X
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
sudo apt -y install git docker.io
git clone https://git.lcomrade.su/root/lenpaste.git
cd ./lenpaste/
git checkout vX.X
sudo docker build -t localhost/lenpaste:latest ./
```

The `localhost/lenpaste:latest` image should now have appeared on your local machine.
You can use it in `docker-compose.yml` or copy it to another machine.



## Other documentation
For all:
- [Frequently Asked Questions](FAQ.md)

For instance administrators:
- [Reverse proxy: Nginx](docs/reverse_proxy_nginx.md)
- [Database: PostgreSQL](docs/db_postgresql.md)
- [Rate limiting](docs/ratelimits.md)
- [Add Lenpaste to Search Engines](docs/search_engines.md)
- [Make Lenpaste server private](docs/private_server.md)
- [Themes for WEB interface](docs/themes.md)

For contributors:
- [Translate on Codeberg Weblate](https://translate.codeberg.org/projects/lenpaste/)
- [Themes for WEB interface](docs/themes.md)

Lenpaste API:
- [Lenpaste API](https://paste.lcomrade.su/docs/apiv1)
- [Libraries for working with API](https://paste.lcomrade.su/docs/api_libs)

Community driven packages:
- [Lenpaste | TrueCharts](https://truecharts.org/docs/charts/incubator/lenpaste/)
- [Lenpaste | Easypanel](https://easypanel.io/docs/templates/lenpaste)

Might be interesting:
- [How to Install LenPaste on Your Synology NAS](https://mariushosting.com/how-to-install-lenpaste-on-your-synology-nas/) (WEB site)
- [Pastebin Clone in Docker with Lenpaste](https://www.youtube.com/watch?v=YxcHxsZHh9A) (YouTube video)



## Contribute
What can I do?
- Translate Lenpaste to you Language: [Codeberg Weblate/Lenpaste](https://translate.codeberg.org/projects/lenpaste/)
- Write an article (DevTo, Medium, your website, and so on) or make a video (YouTube, PeerTube, and so on).
  A link to your article/video will be included in this README.
- Create or update a package:
	- Create NixOS package.
	- Update TrueCharts package.
	- Other.
- Install the Lenpaste server and add it to [Lenmonitor](https://monitor.lcomrade.su/).
- Recommend Lenpaste to your friends.



## Contacts
- Matrix room: [`#lenpaste:lcomrade.su`](https://matrix.to/#/#lenpaste:lcomrade.su)
- Contact me: Leonid Maslakov \<root@lcomrade.su\>



## Donate
All donations will go to Leonid Maslakov, for now the sole developer:
- Qiwi: https://qiwi.com/n/LCOMRADE
- YooMoney (aka Yandex Money): https://yoomoney.ru/to/4100118011659535
