FROM golang:1.23-alpine AS builder

ARG GOARCH
ENV GOARCH=${GOARCH}


# Environment variables for building Go application
RUN echo "Building Go application"
RUN echo "GOARCH: $GOARCH"
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=${GOARCH}


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o smarthome-monitor

FROM alpine:latest
ENV DISCONNECT_TIMEOUT=1000
WORKDIR /root/
COPY --from=builder /app/.env .
COPY --from=builder /app/smarthome-monitor .
#COPY config.yaml .
EXPOSE 2112

# Command to run the application
CMD ["./smarthome-monitor"]
