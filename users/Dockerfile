# Étape 1: Utiliser une image de base officielle de Go
FROM golang:1.23-alpine AS builder

# Définir le répertoire de travail dans le conteneur
WORKDIR /app

# Copier le fichier go.mod et go.sum pour télécharger les dépendances
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier le reste du code source
COPY . .

# Compiler l'application
RUN CGO_ENABLED=0 GOOS=linux go build -o microservice .

# Étape 2: Utiliser une image minimale pour l'exécution
FROM alpine:latest

# Installer les dépendances nécessaires (si nécessaire)
# RUN apk --no-cache add ca-certificates

# Définir le répertoire de travail
WORKDIR /root/

# Copier le binaire compilé depuis l'étape précédente
COPY --from=builder /app/microservice .

# Exposer le port sur lequel le microservice écoute
EXPOSE 8080

# Commande pour exécuter le microservice
CMD ["./microservice"]