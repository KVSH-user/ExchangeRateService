FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0

RUN apk add git && apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
RUN apk add --no-cache ca-certificates tzdata
COPY . .
RUN go build -ldflags="-s -w" -o /build/main cmd/exchangerateservice/main.go

FROM alpine:3.19

ENV TZ=Europe/Moscow

WORKDIR /app
COPY --from=builder /build/main ./main
COPY --from=builder /build/config ./config
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV CONFIG_PATH=/config/docker.yaml
EXPOSE 9049

CMD ["./main"]
