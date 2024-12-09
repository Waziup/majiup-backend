# Stage 1: Build the executable
FROM golang:1.16-alpine AS backend-build
RUN apk add zip
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -a -installsuffix cgo -o majiup .
RUN zip index.zip docker-compose.yml package.json

# Stage 2: Create the final runtime image
FROM alpine:latest
WORKDIR /root/app

COPY --from=backend-build /app/majiup ./
COPY --from=backend-build /app/index.zip /

EXPOSE 8081

CMD ["./majiup"]