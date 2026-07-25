FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /proxydoctor ./cmd/cli
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.20

RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /proxydoctor /usr/local/bin/proxydoctor
COPY --from=builder /server /usr/local/bin/server

EXPOSE 8080

ENTRYPOINT ["proxydoctor"]
