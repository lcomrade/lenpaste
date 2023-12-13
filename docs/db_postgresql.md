# Database: PostgreSQL
Here is an example of the basic Docker Compose config (`docker-compose.yml`):
```yaml
version: "2"

services:
  lenpaste:
    image: ghcr.io/lcomrade/lenpaste:X.X
    restart: always
    environment:
      - LENPASTE_DB_DRIVER=postgres
      - LENPASTE_DB_SOURCE=postgres://lenpaste:pass@postgres/lenpaste?sslmode=disable
    volumes:
      - "${PWD}/data:/data"
      - "/etc/timezone:/etc/timezone:ro"
      - "/etc/localtime:/etc/localtime:ro"
    ports:
      - "80:80"
    depends_on:
      - postgres

  postgres:
    image: docker.io/library/postgres:16.1
    restart: always
    environment:
      - PGDATA=/var/lib/postgresql/data/pgdata
      - POSTGRES_USER=lenpaste
      - POSTGRES_PASSWORD=pass
    volumes:
      - "${PWD}/data/postgres:/var/lib/postgresql/data"
```
