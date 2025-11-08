# Build Stage
# First pull Golang image
FROM golang:1.25.3-alpine as build-env

# Set environment variables
ENV APP_NAME paymentservice

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Set working directory for module-aware build
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies (this layer will be cached if go.mod/go.sum don't change)
RUN go mod download

# Copy application source code
COPY . .

# Build application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /$APP_NAME .

# Run Stage
FROM alpine

# Set envirment variable
ENV APP_NAME paymentservice

# Copy only required data into this image
COPY --from=build-env /$APP_NAME .
COPY ./config/dbconfig.yaml ./config/dbconfig.yaml

# Expose application port
EXPOSE 3000

# Start app
CMD ./$APP_NAME