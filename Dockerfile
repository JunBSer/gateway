FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/gateway ./cmd/gw/main.go


FROM alpine

WORKDIR /app

RUN apk add --no-cache dumb-init

COPY --from=builder /app/gateway /app/gateway
COPY --from=builder /app/configs   /app/configs/
COPY --from=builder /app/docs    /app/docs/

EXPOSE 8080

CMD ["/app/gateway"]
