# Stage 1: Build the executable
FROM golang:1.17 AS builder

WORKDIR /go/src/app

COPY . .


# Stage 2: Create the final runtime image
FROM alpine:latest
WORKDIR /root/app

COPY --from=backend-build /app/majiup ./
COPY --from=backend-build /app/index.zip /

WORKDIR /root/
COPY --from=builder /go/src/app/main .

# Expose the port on which the Go application will run (if needed)
EXPOSE 8082

# Command to run the executable
CMD ["./main"]