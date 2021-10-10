# BUILD
FROM golang:1.15.9-alpine as build

WORKDIR /build

COPY . ./

RUN mkdir ./dist/ && go build -ldflags="-w -s" -o ./dist/lenpaste ./cmd/lenpaste.go


# RUN
FROM alpine:latest as run

WORKDIR /app

COPY --from=build /build/dist/* ./

RUN chmod +x /app/lenpaste

VOLUME /app/data

CMD [ "./lenpaste" ]
