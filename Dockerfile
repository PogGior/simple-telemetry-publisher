# Build stage
FROM golang:1.22.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make compile

# Final stage
FROM scratch
COPY --from=builder /app/simple-telemetry-publisher /app/
WORKDIR /app
CMD ["./simple-telemetry-publisher"]