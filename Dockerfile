# Stage 1: Build the executable
FROM golang:1.17 AS builder

WORKDIR /go/src/app

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Create a minimal image to run the executable
FROM scratch

WORKDIR /root/

# Copy the binary built in Stage 1
COPY --from=builder /go/src/app/main .

# Expose the port on which the Go application will run (if needed)
EXPOSE 8082

# Command to run the executable
CMD ["./main"]