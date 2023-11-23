# Stage 1: Compile go app
FROM golang:1.16-alpine AS backend-build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o majiup .

# Stage 2: build front-end
FROM node:18.13.0-alpine3.16 AS frontend-build
WORKDIR /app
COPY majiup-waziapp/package.json  ./
RUN npm install --force
COPY majiup-waziapp/. .
RUN npm run build

# Stage 3: Create the final runtime image
FROM alpine:latest
WORKDIR /root/app

COPY --from=backend-build /app/majiup ./
COPY --from=frontend-build /app/dist/ ./serve/

EXPOSE 8081

CMD ["./majiup"]
