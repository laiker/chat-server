FROM golang:1.22-alpine AS builder

COPY . /github.com/laiker/chat-server/
WORKDIR /github.com/laiker/chat-server/
RUN go mod download
RUN go build -v -o ./bin/chat-server ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/laiker/chat-server/.env .
COPY --from=builder /github.com/laiker/chat-server/service.pem .
COPY --from=builder /github.com/laiker/chat-server/bin/chat-server .

CMD ["./chat-server"]