# Use a minimal Go base image
FROM golang:1.23 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . /src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/chaos-controller ./cmd/chaos-controller/main.go

# Use a smaller base image to keep the image size smaller
FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/chaos-controller /app/chaos-controller
RUN chmod +x /app/chaos-controller

ENTRYPOINT ["/app/chaos-controller"]
