# Official Go Alpine Base Image for building the application
FROM golang:1.24-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the Go binary
RUN go build -o global-api cmd/global-api/main.go

# Final Image Creation Stage using a lightweight Alpine image
FROM alpine:3.21

# Set the working directory
WORKDIR /root/

# Install any necessary dependencies (e.g., for running Go binaries or for configuration file access)
RUN apk add --no-cache libc6-compat bash

# Copy the built Go binary from the builder image
COPY --from=builder /app/global-api .

# Copy the config file to the container (make sure the path is correct)
COPY config/ /root/config/

# Copy the credentials file to the container (make sure the path is correct)
COPY credentials/ /root/credentials/

# Copy the seed.sql file to the container at the correct path
COPY ./internal/database/seed.sql /app/internal/database/seed.sql

# Set read permissions for the seed.sql file
RUN chmod 644 /app/internal/database/seed.sql

# Copy the wait-for-it.sh script into the container
COPY ./scripts/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Expose the necessary port
EXPOSE 8001

ENV SENBOX_API_KEY=978acdfaa332139818f3467664f7b72c

# Set the entrypoint to wait for MariaDB to be ready before starting the application
CMD /wait-for-it.sh go-main-service-db:3306 -- ./global-api config/config.yaml
