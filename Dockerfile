# Use the official Go image as the base image
FROM golang:1.22.5-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install make and other necessary tools, including gcc and musl-dev for CGO
RUN apk add --no-cache make git nodejs npm gcc musl-dev

# Enable CGO
ENV CGO_ENABLED=1

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy package.json and package-lock.json (if you have them)
COPY package*.json ./

# Install npm dependencies
RUN npm install

# Copy the source code into the container
COPY . .

# Build the application with CGO enabled
RUN make build

# Start a new stage from scratch
FROM alpine:latest

# Install necessary runtime libraries
RUN apk add --no-cache libc6-compat

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy any other necessary files (e.g., templates, static files)
#COPY --from=builder /app/templates ./templates
#COPY --from=builder /app/static ./static

# Create directories for persistent storage
RUN mkdir -p /public /db

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]