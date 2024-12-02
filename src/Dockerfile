FROM golang:1.22

WORKDIR /app

# Copier les fichiers de configuration
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copier le code source
COPY backend/ .

# Copier les fichiers frontend
COPY frontend/ ./frontend/

# Construire l'application
RUN go build -o main .

# Exposer le port
EXPOSE 8080

# Démarrer l'application
CMD ["./main"]
