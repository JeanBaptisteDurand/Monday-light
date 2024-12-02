FROM golang:1.22

WORKDIR /app

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@latest

# Copy configuration files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./
COPY frontend/ ./frontend/

# Expose the application port
EXPOSE 8080

# Run Air for hot reloading
CMD ["air"]
