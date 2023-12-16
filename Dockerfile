# BUILD
FROM docker.io/library/debian:bookworm-20231120-slim as build

WORKDIR /build

RUN sed -i '/^URIs:/d' /etc/apt/sources.list.d/debian.sources && \
    sed -i 's/^# http/URIs: http/' /etc/apt/sources.list.d/debian.sources && \
    apt-get update -o Acquire::Check-Valid-Until=false && \
    apt-get install --no-install-recommends -y make git golang gcc ca-certificates && \
    apt-get clean

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 make && \
    mkdir -p ./dist/data/

# RUN
FROM scratch

COPY --from=build /build/dist/bin/lenpaste /build/dist/data /

VOLUME /data
EXPOSE 80/tcp

ENV LENPASTE_DB_SOURCE=/data/lenpaste.db

ENTRYPOINT [ "/lenpaste" ]
