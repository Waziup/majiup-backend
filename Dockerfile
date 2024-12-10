# Stage 1: Compile go app
FROM golang:1.18-alpine AS backend-build
RUN apk add zip 
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o majiup .
RUN zip index.zip docker-compose.yml package.json


# Stage 2: Create the final runtime image
FROM alpine:latest
WORKDIR /root/app

COPY --from=backend-build /app/majiup ./
COPY --from=backend-build /app/index.zip /

EXPOSE 8081

CMD ["./majiup"]
