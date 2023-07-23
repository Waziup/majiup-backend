# Stage 1: Build the Go application
FROM golang:1.16-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o majiup_backend .


# Stage 2: Create the final runtime image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the previous stage
COPY --from=build /app/majiup_backend ./

# Set the entry point command
CMD ["./majiup_backend"]

