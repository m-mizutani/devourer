# FROM golang:1.21.5-alpine3.18 AS build-go
FROM golang:1.21.5-bullseye AS build-go
ARG BUILD_VERSION
COPY . /app
WORKDIR /app
RUN apt-get update && apt-get install -y libpcap-dev libpcap0.8
RUN go build -o devourer -ldflags "-X github.com/m-mizutani/devourer/pkg/domain/types.AppVersion=${BUILD_VERSION}" .

# FROM gcr.io/distroless/base
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y libpcap0.8
COPY --from=build-go /app/devourer /devourer

ENTRYPOINT ["/devourer"]
