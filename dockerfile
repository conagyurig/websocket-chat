# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o myapp ./cmd/server

# Stage 2: Run the Go binary
FROM alpine:latest

# Set the Current Working Directory inside the container for the runtime image
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/myapp .
COPY --from=build /app/cmd/server/.env .

# Expose port 8000 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]