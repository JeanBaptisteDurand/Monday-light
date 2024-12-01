# Build stage
FROM golang:1.21 as builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./ # Ensure these files exist
RUN go mod download
COPY backend/ . # Copy the backend folder
RUN go build -o main .

# Runtime stage
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .
COPY ./frontend ./frontend # Adjust the path to match the context
EXPOSE 8080
CMD ["./main"]
