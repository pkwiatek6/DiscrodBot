## Build
FROM golang:1.21-alpine AS build
WORKDIR /app

# Copy and download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . ./

# Build the Go application
RUN go build -o /Discrodbot

## Deploy
FROM alpine:latest

WORKDIR /

COPY --from=build /Discrodbot /Discrodbot
CMD [ "/Discrodbot" ]