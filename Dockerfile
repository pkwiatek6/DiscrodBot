## Build
FROM golang:1.21-alpine AS build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /Discrodbot

## Deploy
FROM alpine:latest

WORKDIR /

# Optional: Add netcat to help wait for MongoDB in your app logic
RUN apk add --no-cache netcat-openbsd

COPY --from=build /Discrodbot /Discrodbot
CMD [ "/Discrodbot" ]