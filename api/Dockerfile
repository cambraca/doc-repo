FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -ldflags "-X main.UTCBuildTime=$(TZ=UTC date +%Y-%m-%dT%H:%M:%SZ)" main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]