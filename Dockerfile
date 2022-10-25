# BUILD
FROM golang:1.18.7-alpine as build

WORKDIR /build

RUN apk update && apk upgrade && apk add --no-cache make git gcc musl-dev

COPY ./go.mod ./
COPY go.sum ./
RUN go mod download -x

COPY . ./

RUN make


# RUN
FROM alpine:latest as run

WORKDIR /

COPY --from=build /build/dist/bin/* /usr/local/bin/

COPY ./entrypoint.sh /
RUN chmod 755 /entrypoint.sh && mkdir -p /data/

VOLUME /data
EXPOSE 80/tcp

CMD [ "/entrypoint.sh" ]
