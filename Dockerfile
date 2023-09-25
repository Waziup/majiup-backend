# Stage 1: Build the Go application
FROM golang:1.16-alpine AS build

RUN apk add zip 

# Set the working directory inside the container
WORKDIR /app


# Copy the Go module files
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

RUN zip index.zip docker-compose.yml package.json

# Build the Go application
RUN go build -o majiup .


# Stage 2: Create the final runtime image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/app

# Copy the binary from the previous stage
COPY --from=build /app/majiup ./
COPY --from=build /app/index.zip ./

# Copy the React files from the 'serve' directory in the root directory
COPY serve ./serve

# Expose port 8080
EXPOSE 8081

# Set the entry point command
CMD ["./majiup"]
