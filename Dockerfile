FROM golang:1.22.1-alpine as builder

WORKDIR /app

COPY . .

ARG VERSION

RUN set -ex \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-w -s -X main.Version=${VERSION}" -o estkme-rlpa-server main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/estkme-rlpa-server /app/estkme-rlpa-server

RUN set -ex \
    && apk add --no-cache gcompat ca-certificates pcsc-lite-libs libcurl \
    && update-ca-certificates \
    && chmod +x /app/estkme-rlpa-server

EXPOSE 1888

CMD ["/app/estkme-rlpa-server"]
