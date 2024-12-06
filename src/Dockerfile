FROM golang:1.23.3

WORKDIR /app

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy configuration files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./

# Expose the application port
EXPOSE 8080

# Run Air for hot reloading
CMD ["air"]
