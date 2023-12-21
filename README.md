**Lenpaste** is a web service that allows you to share notes anonymously, an alternative to `pastebin.com`.

> I don't have much time to maintain this project right now, but I didn't request it. Open an issue or write to me by e-mail, I will try to answer.
>
> Also the source code I wrote seems to me now very bad. If I have time I will definitely refactor it:)
>
> \- Leonid Maslakov \<root@lcomrade.su\>


## Features
- No need to register.
- Supports multiple languages.
- Uses cookies only to store settings.
- Can work without JavaScript.
- Has its own API.
- Open source and self-hosted.



## Public servers list
| Server                                             | Description                                                                                                       |
| ---------------------------------------------------| ----------------------------------------------------------------------------------------------------------------- |
| [paste.lcomrade.su](https://paste.lcomrade.su)     | Server managed by the Lenpaste developer.                                                                         |
| [code.dbt3ch.com](https://code.dbt3ch.com)         | Server is managed by DB Tech. He made a [video about Lenpaste v1.1](https://www.youtube.com/watch?v=YxcHxsZHh9A). |
| [notepad.co.il](https://notepad.co.il)             | Server managed by Shlomi Porush. He reported the bug and made some suggestions.                                   |
| [lenp.pardesicat.xyz](https://lenp.pardesicat.xyz) | Server managed by Pardesi_Cat. He translated Lenpaste into Bengali and helped correct the documentation.          |



## Launch your own server
### Simple setup
1. If you don't already have Docker and Docker Compose installed, do so:
```
apt-get install -y docker.io docker-compose
```

2. Use a file like this `docker-compose.yml`:
```yaml
version: "3.4"

services:
  lenpaste:
    # There are images for x64, ARM64, ARM v7, ARM v6.
    # The Raspberry Pi is supported, including the latest 64-bit versions.
    # Replace "X.X" to target Lenpaste version.
    image: ghcr.io/lcomrade/lenpaste:X.X
    volumes:
      - "${PWD}/data:/data"
    ports:
      - "80:80"
```

3. Execute while in the directory where `docker-compose.yml` is located:
```
docker-compose pull && docker-compose up -d
```

If you want to install updates, run: `docker-compose pull && docker-compose up -d && docker system prune -a -f`


### Lenpaste configuration
#### HTTP
The `LENPASTE_ADDRESS` environment variable specifies the `ADDRESS:PORT` at which Lenpaste will expect HTTP connections.
The default is `:80`.


#### Database
The `LENPASTE_DB_DRIVER` environment variable specifies the database to be used.
The default is `sqlite3`, possible values are `sqlite3` and `postgres`.

The `LENPASTE_DB_SOURCE` environment variable specifies the data to connect to the database.
In case of SQLite3 the default value is `/data/lenpaste.db`, for other databases it is necessary to specify the value explicitly.

The `LENPASTE_DB_MAX_OPEN_CONNS` environment variable specifies the maximum number of database connections that Lenpaste can open at one time.
The default is `25`.

The `LENPASTE_DB_MAX_IDLE_CONNS` environment variable specifies the maximum number of database connections that Lenpaste can have in idle state.
The default is `5`.

The `LENPASTE_DB_CLEANUP_PERIOD` environment variable specifies the time after which Lenpaste should delete expired pastes.
If the user tries to open an expired paste that has not yet been cleaned, the user will receive a 404 error.
The default is `1m` (1 minute).


#### Search engines
The `LENPASTE_ROBOTS_DISALLOW` environment variable prohibits or allows search engine robots (such as Google) to index your Lenpaste instance via the robots.txt file.
The default is `false` (allow indexing).


#### Storage limits
The `LENPASTE_TITLE_TITLE_MAX_LENGTH` environment variable sets the maximum size of the paste header.
The default is `100`.

The `LENPASTE_BODY_MAX_LENGTH` environment variable sets the maximum size of the paste body.
The default is `20000`.

The `LENPASTE_MAX_PASTE_LIFETIME` environment variable sets the maximum paste life time.
Examples of `10m`, `1h 30m`, `12h`, `7w`, `30d`, `365d` values. The default is `unlimited`.


#### Rate limits
The environment variables `LENPASTE_GET_PASTES_PER_5MIN`, `LENPASTE_GET_PASTES_PER_15MIN`, `LENPASTE_GET_PASTES_PER_1HOUR`
set the maximum number of pastes that can be VIEWED in 5, 15 or 60 minutes from one IP.
The default values is `50`, `100`, `500`.

The environment variables `LENPASTE_NEW_PASTES_PER_5MIN`, `LENPASTE_NEW_PASTES_PER_15MIN`, `LENPASTE_NEW_PASTES_PER_1HOUR`
set the maximum number of pastes that can be CREATED in 5, 15 or 60 minutes from one IP.
The default values is `15`, `30`, `40`.

To turn off any rait limit just set it to `0`.
Lenpaste remembers the limit only until the restart, after the restart the limit count starts again.


#### Access control
If the file `/data/lenpasswd` is present, the server will prompt for a login and password to create the paste.
The file format is `LOGIN:PLAIN_PASSWORD` on each line.


#### Information about server
The `LENPASTE_ADMIN_NAME` environment variable sets the name of the server administrator.

The `LENPASTE_ADMIN_MAIL` environment variable sets the email of the server administrator.

If the `/data/about` text file exists, its contents will be displayed at the top of the "About" page.

The `/data/rules` text file may contain human-readable rules for using the server (it is not a legal document).

The `/data/terms` text file may contain "Terms of Use" written in "legal language".


#### Web interface
The `LENPASTE_UI_DEFAULT_LIFETIME` environment variable sets the default paste lifetime selected in the WEB interface.
Examples of values of `10min`, `1h`, `1d`, `2w`, `6mon`, `1y`.
The first available value in the list will be selected by default.

The `LENPASTE_UI_DEFAULT_THEME` environment variable sets the default theme to be used in the WEB interface.
The default is `dark`.

In the `/data/themes/` directory, the administrator can place custom themes for WEB interface.
You can create a custom theme for the WEB interface based on the themes located in `./internal/web/data/theme/`.


### Lenpaste + Nginx
Here is an example of the basic Nginx config (`/etc/nginx/nginx.conf`) that will work with Lenpaste:
```nginx
events {
	worker_connections 1024;
}

http {
	error_log /var/log/nginx/error.log warn;
	server_tokens off; # Disables emitting nginx version on error pages and in the “Server” response header field.

	ssl_protocols TLSv1.2 TLSv1.3; # TLSv1.2 enables HTTPS on older devices.

	client_max_body_size 1M;
	client_body_timeout 300s;

	proxy_http_version 1.1;
	
	# Required for Lenpaste to work correctly.
	proxy_set_header Host $host;
	proxy_set_header X-Real-IP $remote_addr;
	proxy_set_header X-Forwarded-For $remote_addr;
	proxy_set_header X-Forwarded-Proto $scheme;
}

# HTTP
server {
	server_name YOUR_DOMAIN;
	listen 80;
	listen [::]:80;

	access_log /var/log/nginx/YOUR_DOMAIN.access.log combined;

	location / {
		proxy_pass http://localhost:8000/;
		#return 301 https://$host$request_uri; - redirect to HTTPS
	}

	# Required for Lets Encrypt
	location /.well-known/acme-challenge/ {
		root /var/www/letsencrypt/;
	}
}

# HTTPS
server {
	server_name YOUR_DOMAIN;
	listen 443 ssl http2;
	listen [::]:443 ssl http2;
	ssl_certificate /etc/letsencrypt/live/YOUR_DOMAIN/fullchain.pem;
	ssl_certificate_key /etc/letsencrypt/live/YOUR_DOMAIN/privkey.pem;

	access_log /var/log/nginx/YOUR_DOMAIN.access.log combined;
	
	location / {
		proxy_pass http://localhost:8000/;
	}
}
```

### Lenpaste + Postgres
Here is an example of the basic Docker Compose config (`docker-compose.yml`) that will make Lenpaste work with Postgres:
```yaml
version: "3.4"

services:
  postgres:
    image: docker.io/library/postgres:16.1
    environment:
      - PGDATA=/var/lib/postgresql/data/pgdata
      - POSTGRES_USER=lenpaste
      - POSTGRES_PASSWORD=pass
    volumes:
      - "${PWD}/data/postgres:/var/lib/postgresql/data"

  lenpaste:
    image: ghcr.io/lcomrade/lenpaste:X.X
    restart: on-failure:10
    environment:
      - LENPASTE_DB_DRIVER=postgres
      - LENPASTE_DB_SOURCE=postgres://lenpaste:pass@postgres/lenpaste?sslmode=disable
    volumes:
      - "${PWD}/data:/data"
    ports:
      - "80:80"
    depends_on:
      - postgres
```



## Build from source code
### Build Docker image (recommended)
**Why is it necessary?**
An official image may not support your architecture e.g. MIPS, PowerPC, etc.
So you can build your own image to run on an officially unsupported architecture
(of course you have to rebuild it every time you update Lenpaste).

On Debian/Ubuntu:
```bash
export LENPASTE_VERSION=X.X
sudo apt -y install wget docker.io
wget -O ./lenpaste-$LENPASTE_VERSION.tar.gz https://github.com/lcomrade/lenpaste/releases/download/v$LENPASTE_VERSION/lenpaste-$LENPASTE_VERSION.tar.gz
tar -xf ./lenpaste-$LENPASTE_VERSION.tar.gz
cd ./lenpaste-$LENPASTE_VERSION/
sudo docker build -t localhost/lenpaste:$LENPASTE_VERSION ./
```

The `localhost/lenpaste:X.X` image should now have appeared on your local machine.
You can use it in `docker-compose.yml` or copy it to another machine.


### Build binary
On Debian/Ubuntu:
```bash
export LENPASTE_VERSION=X.X
sudo apt -y install wget make golang gcc libc6-dev
wget -O ./lenpaste-$LENPASTE_VERSION.tar.gz https://github.com/lcomrade/lenpaste/releases/download/v$LENPASTE_VERSION/lenpaste-$LENPASTE_VERSION.tar.gz
tar -xf ./lenpaste-$LENPASTE_VERSION.tar.gz
cd ./lenpaste-$LENPASTE_VERSION/
make
```

You can find the result of the build in the `./dist/` directory.



## Other documentation
Read more about [Lenpaste API](https://paste.lcomrade.su/docs/apiv1).

Might be interesting:
- [How to Install LenPaste on Your Synology NAS](https://mariushosting.com/how-to-install-lenpaste-on-your-synology-nas/) (WEB site)
- [Lenpaste | TrueCharts](https://truecharts.org/docs/charts/incubator/lenpaste/) (WEB site)
- [Pastebin Clone in Docker with Lenpaste](https://www.youtube.com/watch?v=YxcHxsZHh9A) (YouTube video)



## Contribute
What can I do?
- Write an article (DevTo, Medium, your website, and so on) or make a video (YouTube, PeerTube, and so on).
  A link to your article/video will be included in this README.
- Create or update a package:
	- Create NixOS package.
	- Update TrueCharts package.
	- Other.
- Recommend Lenpaste to your friends.



## Contacts
- Matrix room: [`#lenpaste:lcomrade.su`](https://matrix.to/#/#lenpaste:lcomrade.su)
- Contact me: Leonid Maslakov \<root@lcomrade.su\>
