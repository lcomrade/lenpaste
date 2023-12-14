# BUILD
FROM docker.io/library/debian:bookworm-20231120-slim as build

WORKDIR /build

RUN sed -i '/^URIs:/d' /etc/apt/sources.list.d/debian.sources && \
    sed -i 's/^# http/URIs: http/' /etc/apt/sources.list.d/debian.sources && \
    apt-get update -o Acquire::Check-Valid-Until=false && \
    apt-get install -y make git golang gcc && \
    apt-get clean

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 make


# RUN
FROM scratch

COPY --from=build /build/dist/bin/lenpaste /lenpaste

VOLUME /data
EXPOSE 80/tcp

CMD [ "/lenpaste" ]
