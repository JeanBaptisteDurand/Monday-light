# Build stage
FROM golang:1.21 as builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN go build -o main .

# Runtime stage
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .
COPY frontend ./frontend
EXPOSE 8080
CMD ["./main"]
