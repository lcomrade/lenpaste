# BUILD
FROM golang:1.18.1-alpine as build

WORKDIR /build

RUN apk update && apk upgrade && apk add --no-cache make gcc musl-dev

COPY . ./

RUN make


# RUN
FROM alpine:latest as run

WORKDIR /app

COPY --from=build /build/dist/* ./

EXPOSE 8000/tcp

CMD [ "./lenpaste" ]
