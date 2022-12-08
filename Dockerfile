# BUILD
FROM golang:1.18.9-alpine3.17 as build

WORKDIR /build

RUN apk add --no-cache make=4.3-r1 git=2.38.1-r0 gcc=12.2.1_git20220924-r4 musl-dev=1.2.3-r4

COPY ./go.mod ./
COPY go.sum ./
RUN go mod download -x

COPY . ./

RUN make


# RUN
FROM alpine:3.17.0 as run

WORKDIR /

COPY --from=build /build/dist/bin/* /usr/local/bin/

COPY ./entrypoint.sh /
RUN chmod 755 /entrypoint.sh && mkdir -p /data/

VOLUME /data
EXPOSE 80/tcp

CMD [ "/entrypoint.sh" ]
