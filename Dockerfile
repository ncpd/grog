ARG arch=amd64

# Base
FROM golang:1.15-alpine as base
WORKDIR /build
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

# ARM v7 builder
FROM base as armv7
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build

# AMD64 builder
FROM base as amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# COPY --from doesn't support args
FROM ${arch} as build

# Production
FROM scratch
WORKDIR /app
COPY --from=build /build/grog .
CMD [ "/app/grog" ]