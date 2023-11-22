# Stage 1: Compile go app
FROM golang:1.16-alpine AS backend-build
RUN apk add zip 
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o majiup .
RUN zip index.zip docker-compose.yml package.json

# Stage 2: build front-end
FROM node:18.13.0-alpine3.16 AS frontend-build
WORKDIR /app
COPY majiup-waziapp/package.json  ./
RUN npm install --force
COPY majiup-waziapp/. .
#ENV VITE_BACKEND_URL="http://192.168.0.112:8081"
RUN npm run build
RUN ls dist/
# Stage 3: Create the final runtime image
FROM alpine:latest
WORKDIR /root/app

COPY --from=backend-build /app/majiup ./
COPY --from=backend-build /app/index.zip /
COPY --from=frontend-build /app/dist/* ./serve/

EXPOSE 8081

CMD ["./majiup"]
