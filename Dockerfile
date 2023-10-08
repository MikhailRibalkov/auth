FROM golang:1.21.1-alpine AS builder

COPY . /github.com/MikhailRibalkov/auth/source/
WORKDIR /github.com/MikhailRibalkov/auth/source/

RUN go mod download
RUN go build -o ./bin/auth-server cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/MikhailRibalkov/auth/source/bin/auth-server .

CMD ["./auth-server"]