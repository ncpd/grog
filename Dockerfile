FROM golang:1.15-alpine as builder

WORKDIR /build

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch

COPY --from=builder /build/grog .

CMD [ "/grog" ]