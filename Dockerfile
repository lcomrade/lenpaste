# BUILD
FROM docker.io/library/debian:bookworm-20231120-slim as build

WORKDIR /build

RUN sed -i '/^URIs:/d' /etc/apt/sources.list.d/debian.sources && \
    sed -i 's/^# http/URIs: http/' /etc/apt/sources.list.d/debian.sources && \
    apt-get update -o Acquire::Check-Valid-Until=false && \
    apt-get install --no-install-recommends -y make git golang gcc libc6-dev ca-certificates && \
    apt-get clean

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . ./
RUN make

# RUN
FROM docker.io/library/debian:bookworm-20231120-slim

COPY ./entrypoint.sh /

RUN mkdir /data/ && chmod +x /entrypoint.sh

VOLUME /data
EXPOSE 80/tcp

COPY --from=build /build/dist/bin/* /usr/local/bin/

CMD [ "/entrypoint.sh" ]
