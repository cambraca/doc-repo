FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
#RUN chmod +x /app/main  # Explicitly set execute permissions
EXPOSE 8080
CMD ["./main"]